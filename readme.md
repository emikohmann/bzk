## Usage


`go run main.go`

**POST** `/run`


```json
{
    "target": {
        "protocol": "https",
        "base_url": "api.domain.com",
        "method": "GET",
        "headers": {
            "X-Caller-Scopes": "admin"
        },
        "endpoint": {
            "format": "/names/%s",
            "context": [
                [
                    "john",
                    "william",
                    "alex"
                ]
            ]
        }
    },
    "steps": [
        {
            "rpm_from": 1,
            "rpm_to": 1000,
            "duration_seconds": 10
        },
        {
            "rpm_from": 2000,
            "rpm_to": 2000,
            "duration_seconds": 10
        },
        {
            "rpm_from": 1000,
            "rpm_to": 1,
            "duration_seconds": 10
        }
    ]
}
```

## Example output

```json
{
    "step_results": [
        {
            "call_count ": 72,
            "status ": {
                "200": 72
            },
            "duration_seconds ": 10.067134003
        },
        {
            "call_count ": 327,
            "status ": {
                "200": 327
            },
            "duration_seconds ": 10.020911311
        },
        {
            "call_count ": 88,
            "status ": {
                "200": 88
            },
            "duration_seconds ": 10.190898304
        }
    ],
    "duration_seconds": 30.547391958
}
```