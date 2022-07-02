package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type RouterConfig struct {
	Host string
	Port int
}

type Route struct {
	Path   string
	Method string
	handle httprouter.Handle
}

type Handler interface {
	GetRoutes() []Route
}

func newRouter(handlers ...Handler) http.Handler {
	router := httprouter.New()
	for _, h := range handlers {
		for _, r := range h.GetRoutes() {
			router.Handle(r.Method, r.Path, r.handle)
		}
	}
	return router
}

type QueryParameters struct {
	values url.Values
}

func (qp QueryParameters) GetString(key string) string {
	return strings.TrimSpace(qp.values.Get(key))
}

func (qp QueryParameters) GetStringSlice(key string) []string {
	var ss []string
	if qp.values == nil {
		return ss
	}
	for _, v := range qp.values[key] {
		s := strings.TrimSpace(v)
		if len(s) > 0 {
			ss = append(ss, s)
		}
	}
	return ss
}

func (qp QueryParameters) GetInt(key string) (int, error) {
	s := qp.GetString(key)
	if len(s) == 0 {
		return 0, nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("QueryParameters.GetInt: %w", err)
	}
	return i, nil
}

func (qp QueryParameters) GetTime(key string) (time.Time, error) {
	layout := "2006-01-02T15:04"

	s := qp.GetString(key)
	if len(s) == 0 {
		return time.Time{}, nil
	}
	t, err := time.Parse(layout, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("QueryParameters.GetTime: %q does not match time layout %q", s, layout)
	}
	return t, nil
}

func (qp QueryParameters) GetBool(key string, value bool) (bool, error) {
	switch qp.GetString(key) {
	case "":
		return value, nil
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, errors.New("QueryParameters.GetBool: invalid format: expect \"true\" or \"false\"")
	}
}
