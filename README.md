# rest-calculator-service

## Запуск

```
cat ./computing-core/.env.example > ./computing-core/.env && \
cat ./gateway/.env.example > ./gateway/.env && \
export PROJECT_ROOT=$(pwd) && \
docker compose up --build
```

## Примеры запросов

С примерами запросов можно ознакомиться [здесь](/example-requests.md).