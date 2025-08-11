# HTTP-сервер «Календарь» (WB L2.18)

#####  HTTP-сервер для работы с небольшим календарем событий.

---

Для работы с проектом добавьте в корневую директорию файл `.env` со следующими полями:

```
CONFIG_PATH=./config.yaml
HTTP_ADDRESS=:8123
HTTP_TIMEOUT=4s
HTTP_IDLE_TIMEOUT=60s
HTTP_USER=admin
HTTP_PASSWORD=secret
```

### Запуск сервера

Для запуска сервера необходимо выполнить две команды:
1. ```docker-compose build``` (перед сборкой выключите VPN)
2. ```docker-compose up```

После сборки сервер будет доступен на порту `:8777`

### CRUD операции:

#### Создание события
```http
POST http://localhost:8777/create_event
Content-Type: application/json

{
    "id": 1,
    "user_id": 1,
    "date": "2025-08-11",
    "title": "but potatkkko"
}
```

#### Обновление события
```http
POST http://localhost:8777/update_event
Content-Type: application/json

{
  "id": 1,
  "date": "2023-12-26",
  "title": "Boxing Day"
}
```

#### Удаление события
```http
POST http://localhost:8777/delete_event
Content-Type: application/json

{
  "id": 1
}
```

### Получение событий (GET запросы)

#### События на день
```http
GET http://localhost:8777/events_for_day
Content-Type: application/json

{
    "user_id": 1,
    "date": "2025-08-11"
}
```

#### События на неделю
```http
GET http://localhost:8777/events_for_week
Content-Type: application/json

{
    "user_id": 1,
    "date": "2025-08-11"
}
```

#### События на месяц
```http
GET http://localhost:8777/events_for_month
Content-Type: application/json

{
    "user_id": 1,
    "date": "2025-08-11"
}
```
--- 

### Тесты
Для запуска тестов воспользуйтесь командой `make test`