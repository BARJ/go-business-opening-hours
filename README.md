# Go Business Opening Hours

## Starting the application

```bash
sudo docker compose up --build -d
```

## Testing the application

1. Install [Postman](https://www.postman.com/).
2. Import the Postman Collection `Go Business Opening Hours.postman_collection.json`.
3. Execute the HTTP request `List Businesses` and play with the various pre-defined URL Query Parameters for filtering (see "Params" tab.).

## Example

Request:
```bash
GET http://localhost:8080/businesses?open=true&local_time=2022-07-09T00:15
```

Response:
```json
{
    "Businesses": [
        {
            "BusinessID": 1,
            "Name": "Convenience Store",
            "OpeningHours": [
                {
                    "Day": "Monday",
                    "Opens": "00:00",
                    "Closes": "23:59"
                },
                {
                    "Day": "Tuesday",
                    "Opens": "00:00",
                    "Closes": "23:59"
                },
                {
                    "Day": "Wednesday",
                    "Opens": "00:00",
                    "Closes": "23:59"
                },
                {
                    "Day": "Thursday",
                    "Opens": "00:00",
                    "Closes": "23:59"
                },
                {
                    "Day": "Friday",
                    "Opens": "00:00",
                    "Closes": "23:59"
                },
                {
                    "Day": "Saturday",
                    "Opens": "00:00",
                    "Closes": "23:59"
                },
                {
                    "Day": "Sunday",
                    "Opens": "00:00",
                    "Closes": "23:59"
                }
            ]
        },
        {
            "BusinessID": 2,
            "Name": "Late Night Bar",
            "OpeningHours": [
                {
                    "Day": "Monday",
                    "Opens": "17:00",
                    "Closes": "02:00"
                },
                {
                    "Day": "Tuesday",
                    "Opens": "17:00",
                    "Closes": "02:00"
                },
                {
                    "Day": "Wednesday",
                    "Opens": "17:00",
                    "Closes": "02:00"
                },
                {
                    "Day": "Thursday",
                    "Opens": "17:00",
                    "Closes": "02:00"
                },
                {
                    "Day": "Friday",
                    "Opens": "17:00",
                    "Closes": "02:00"
                },
                {
                    "Day": "Saturday",
                    "Opens": "14:00",
                    "Closes": "02:00"
                },
                {
                    "Day": "Sunday",
                    "Opens": "14:00",
                    "Closes": "02:00"
                }
            ]
        },
        {
            "BusinessID": 3,
            "Name": "Restaurant",
            "OpeningHours": [
                {
                    "Day": "Monday",
                    "Opens": "00:00",
                    "Closes": "00:00"
                },
                {
                    "Day": "Tuesday",
                    "Opens": "12:00",
                    "Closes": "15:00"
                },
                {
                    "Day": "Tuesday",
                    "Opens": "18:00",
                    "Closes": "23:00"
                },
                {
                    "Day": "Wednesday",
                    "Opens": "12:00",
                    "Closes": "15:00"
                },
                {
                    "Day": "Wednesday",
                    "Opens": "18:00",
                    "Closes": "23:00"
                },
                {
                    "Day": "Thursday",
                    "Opens": "12:00",
                    "Closes": "15:00"
                },
                {
                    "Day": "Thursday",
                    "Opens": "18:00",
                    "Closes": "23:00"
                },
                {
                    "Day": "Friday",
                    "Opens": "12:00",
                    "Closes": "15:00"
                },
                {
                    "Day": "Friday",
                    "Opens": "18:00",
                    "Closes": "00:30"
                },
                {
                    "Day": "Saturday",
                    "Opens": "12:00",
                    "Closes": "15:00"
                },
                {
                    "Day": "Saturday",
                    "Opens": "18:00",
                    "Closes": "00:00"
                },
                {
                    "Day": "Sunday",
                    "Opens": "12:00",
                    "Closes": "15:00"
                },
                {
                    "Day": "Sunday",
                    "Opens": "18:00",
                    "Closes": "23:00"
                }
            ]
        }
    ]
}
```
