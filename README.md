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

## Структура проекта

```
├───computing-core                       # сервис-калькулятор 
│   ├───cmd
│   │   └───app                          # точка входа в приложение
│   ├───internal
│   │   ├───core                         # общие компоненты  
│   │   │   ├───domain
│   │   │   ├───errors
│   │   │   ├───logger                   # конфигурация и инициализация логгера
│   │   │   └───transport
│   │   │       └───http
│   │   │           ├───middleware
│   │   │           ├───request
│   │   │           ├───response         
│   │   │           └───server           # конфигурация и запуск HTTP-сервера
│   │   ├───features
│   │   │   └───calculator               # фича калькулятора
│   │   │       ├───repository           # слой доступа к данным сервиса-калькулятора (пустой)
│   │   │       ├───service              # слой бизнес-логики сервиса-калькулятора
│   │   │       └───transport            # обработчики HTTP-эндпоинтов сервиса-калькулятора
│   │   │           └───http
│   │   └───utils                        # вспомогательные утилиты
│   │       ├───hash
│   │       ├───http
│   │       ├───jwt
│   │       └───time
│   └───out                              # директория с логами
│       └───logs
└───gateway                              # шлюз-агрегатор
    ├───cmd
    │   └───app                          # точка входа в приложение
    ├───internal
    │   ├───core                         # общие компоненты
    │   │   ├───domain
    │   │   ├───errors
    │   │   ├───logger                   # конфигурация и инициализация логгера
    │   │   └───transport
    │   │       └───http
    │   │           ├───middleware
    │   │           ├───request
    │   │           ├───response
    │   │           └───server           # конфигурация и запуск HTTP-сервера
    │   ├───features
    │   │   └───calculator               # фича калькулятора
    │   │       ├───repository           # слой доступа к данным шлюза-агрегатора (пустой)
    │   │       ├───service              # слой бизнес-логики шлюза-агрегатора
    │   │       └───transport            # обработчики HTTP-эндпоинтов шлюза-агрегатора
    │   │           └───http
    │   └───utils                        # вспомогательные утилиты
    │       ├───hash
    │       ├───http
    │       ├───jwt
    │       └───time
    └───out
        └───logs                         # директория с логами
```