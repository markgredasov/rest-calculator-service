# Примеры запросов

### Успешный запрос

`POST http://localhost:8000/api/v1/calculate`

```
{
    "numbers": [1, 2, 3],
    "operation": "average",
    "params": {
        "power": 3
    }
}
```

Ответ:

```
STATUS: 200 OK
{
    "status": "success",
    "original_numbers": [
        1,
        2,
        3
    ],
    "transformed_numbers": [
        1,
        8,
        27
    ],
    "aggregated_result": 144
}
```

### Запрос с некорректными аргументами

`POST http://localhost:8000/api/v1/calculate`

```
{
    "numbers": [1, 2, 3],
    "operation": "avg",
    "params": {
        "power": 3
    }
}
```

Ответ:
```
STATUS: 400 Bad Request
{
    "error": "invalid argument",
    "message": "failed to decode and validate request"
}
```

### Запрос при недоступности сервиса

`POST http://localhost:8000/api/v1/calculate`

```
{
    "numbers": [1, 2, 3],
    "operation": "average",
    "params": {
        "power": 3
    }
}
```

Ответ:

```
STATUS: 503 Service Unavailable
{
    "error": "service unavailable",
    "message": "failed to calculate result"
}
```