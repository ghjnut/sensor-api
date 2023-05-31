## POST /ingest
```
curl --location --request POST 'http://localhost:8000/ingest' \
--header 'Content-Type: application/json' \
--data-raw '{
    "data": [
        "YYARKx|2022-03-22T21:42:02.362Z|-8",
        "YYARKx|2022-03-22T21:42:04.372Z|-1",
        "YYARKx|2022-03-22T21:42:50.572Z|7"
    ]
}'
```

## GET /device/{id}
```
curl --location --request GET 'http://localhost:8000/device/YYARKx'
```
