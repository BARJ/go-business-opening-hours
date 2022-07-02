package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"regexp"
	"time"
)

type Business struct {
	BusinessID   int
	Name         string
	OpeningHours OpeningHours
}

type OpeningHours struct {
	Periods []OpeningHoursPeriod
}

func (oh OpeningHours) MarshalJSON() ([]byte, error) {
	return json.Marshal(oh.Periods)
}

func (oh *OpeningHours) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &oh.Periods); err != nil {
		return fmt.Errorf("OpeningHours.UnmarshalJSON: %w", err)
	}
	return nil
}

type OpeningHoursPeriod struct {
	Day    Day
	Opens  Clock
	Closes Clock
}

type Day int

const (
	DayUndefined Day = iota
	DayMondy
	DayTuesday
	DayWednesday
	DayThursday
	DayFriday
	DaySaturday
	DaySunday
)

var dayNames = [...]string{"Undefined", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

func (d Day) String() string {
	return dayNames[d]
}

func (d Day) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}
func (d Day) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

var clockPattern = regexp.MustCompile(`^([0-9]|0[0-9]|1[0-9]|2[0-3]):[0-5][0-9](:[0-5][0-9])?$`)

type Clock struct {
	Hours   int
	Minutes int
}

func ParseClock(s string) (*Clock, error) {
	fail := func() (*Clock, error) {
		return nil, fmt.Errorf("Clock.ParseClock: %q does not match clock pattern \"HH:MM[:SS]\"", s)
	}
	if matched := clockPattern.MatchString(s); !matched {
		return fail()
	}
	var err error
	var t time.Time
	for _, layout := range []string{"15:04:05", "15:04"} {
		if t, err = time.Parse(layout, s); err == nil {
			break
		}
	}
	if err != nil {
		return fail()
	}
	return &Clock{Hours: t.Hour(), Minutes: t.Minute()}, nil
}

func (c Clock) String() string {
	return c.time().Format("15:04")
}

func (c Clock) time() time.Time {
	return time.Date(0, 1, 1, c.Hours, c.Minutes, 0, 0, time.UTC)
}

func (c Clock) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Clock) UnmarshalJSON(data []byte) error {
	fail := func(err error) error {
		return fmt.Errorf("Clock.UnmarshalJSON: %w", err)
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fail(fmt.Errorf("failed to unmarshal json into string: %w", err))
	}

	clock, err := ParseClock(s)
	if err != nil {
		return fail(fmt.Errorf("failed to parse clock %w", err))
	}
	*c = *clock
	return nil
}

func (b *Business) Scan(src interface{}) error {
	fail := func(err error) error {
		return fmt.Errorf("Business.Scan: %w", err)
	}

	data, ok := src.([]byte)
	if !ok {
		return fail(errors.New("type assertions to []byte failed"))
	}

	if err := json.Unmarshal(data, &b); err != nil {
		return fail(fmt.Errorf("json serialisation failed: %w", err))
	}
	return nil
}

type BusinessFilter struct {
	Open      bool
	LocalTime time.Time
}

func (bf BusinessFilter) Validate() error {
	if bf.Open && bf.LocalTime.Equal(time.Time{}) {
		return errors.New("BusinessFilter.Validate: require \"LocalTime\"")
	}
	return nil
}

type BusinessHandler struct {
	BusinessStore BusinessStore
}

func (h BusinessHandler) GetRoutes() []Route {
	return []Route{
		{"/businesses", http.MethodGet, h.ListBusinesses},
	}
}

func (h BusinessHandler) ListBusinesses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fail := func(err error, code int) {
		log.Printf("BusinessHandler.ListBusinesses: %v\n", err)
		http.Error(w, http.StatusText(code), code)
	}

	// Validate and parse URL query parameters.
	queryParams := QueryParameters{r.URL.Query()}
	open, err := queryParams.GetBool("open", false)
	if err != nil {
		fail(err, http.StatusBadRequest)
		return
	}
	localTime, err := queryParams.GetTime("local_time")
	if err != nil {
		fail(err, http.StatusBadRequest)
		return
	}
	businessFilter := BusinessFilter{open, localTime}
	if err := businessFilter.Validate(); err != nil {
		fail(err, http.StatusBadRequest)
		return
	}

	// Get businesses from storage.
	businesses, err := h.BusinessStore.ListBusinesses(businessFilter)
	if err != nil {
		fail(err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Businesses []Business
	}{businesses})
}

type BusinessStore struct {
	DB *sql.DB
}

func (s BusinessStore) ListBusinesses(filter BusinessFilter) ([]Business, error) {
	fail := func(err error) ([]Business, error) {
		return nil, fmt.Errorf("BusinessStore.ListBusinesses: %w", err)
	}

	query := `
SELECT JSON_BUILD_OBJECT(
               'BusinessID', b.business_id,
               'Name', b.name,
               'OpeningHours',
               COALESCE((SELECT JSON_AGG(JSON_BUILD_OBJECT('day', oh.day, 'opens', oh.opens, 'closes', oh.closes))
                         FROM opening_hours oh
                         WHERE oh.business_id = b.business_id), JSON_BUILD_ARRAY())
           )
FROM business b
WHERE COALESCE($1, FALSE) IS FALSE
   OR EXISTS(
        SELECT *
        FROM opening_hours oh
        WHERE oh.business_id = b.business_id
          AND oh.opens != oh.closes
          AND ((oh.day = EXTRACT(ISODOW FROM $2::TIMESTAMP) AND
                ((oh.opens = '00:00'::TIME AND oh.closes = '23:59'::TIME) OR
                 ($2::TIME >= oh.opens AND (oh.opens > oh.closes OR $2::TIME <= oh.closes))))
            OR (oh.day + 1 % 7 = EXTRACT(ISODOW FROM $2::TIMESTAMP) AND oh.opens > oh.closes AND
                $2::TIME <= oh.closes)))
`
	rows, err := s.DB.Query(query, filter.Open, filter.LocalTime)
	if err != nil {
		return fail(err)
	}

	businesses := []Business{}
	for rows.Next() {
		var business Business
		if err := rows.Scan(&business); err != nil {
			return fail(err)
		}
		businesses = append(businesses, business)
	}

	if err := rows.Close(); err != nil {
		return fail(err)
	}

	return businesses, nil
}
