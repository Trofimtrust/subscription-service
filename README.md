# Subscription Service

REST API для управления онлайн-подписками пользователей.

## Стек

- Go
- Chi Router
- PostgreSQL
- Docker Compose
- pgx
- Swagger

## Запуск

### 1. Запустить PostgreSQL

```bash
docker compose up -d
```

### 2. Применить миграции

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/subscriptions?sslmode=disable" up
```

### 3. Запустить приложение

```bash
go run main.go
```

## API

### Создать подписку

POST

```
/subscriptions
```

```json
{
    "service_name": "Yandex Music",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
}
```

### Получить список

GET

```
/subscriptions
```

### Получить по ID

GET

```
/subscriptions/{id}
```

### Обновить

PUT

```
/subscriptions/{id}
```

### Удалить

DELETE

```
/subscriptions/{id}
```

### Подсчитать стоимость

GET

```
/subscriptions/cost?user_id=<uuid>&service_name=Yandex%20Music&from=2025-01&to=2025-12
```

## Конфигурация

Используется файл `.env`

```
DB_HOST=localhost
DB_PORT=5432
DB_NAME=subscriptions
DB_USER=postgres
DB_PASSWORD=postgres
```