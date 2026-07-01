# Subscription Service

REST API для управления онлайн-подписками пользователей.

## Стек

- Go
- Chi Router
- PostgreSQL
- Docker Compose
- pgx
- Swagger

---

## Запуск

### 1. Клонировать репозиторий

```bash
git clone https://github.com/Trofimtrust/subscription-service.git
cd subscription-service
```

### 2. Запустить проект

```bash
docker compose up --build
```

После запуска сервис будет доступен по адресу:

```
http://localhost:8080
```

Swagger:

```
http://localhost:8080/swagger/index.html
```

> При первом запуске PostgreSQL автоматически создаст базу данных и выполнит SQL-миграцию.

---

## API

### Создать подписку

**POST**

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

### Получить все подписки

**GET**

```
/subscriptions
```

### Получить подписку по ID

**GET**

```
/subscriptions/{id}
```

### Обновить подписку

**PUT**

```
/subscriptions/{id}
```

### Удалить подписку

**DELETE**

```
/subscriptions/{id}
```

### Рассчитать стоимость подписок

**GET**

```
/subscriptions/cost?user_id=<uuid>&service_name=Yandex%20Music&from=2025-01&to=2025-12
```

---

## Конфигурация

Используются следующие переменные окружения:

```
DB_HOST
DB_PORT
DB_NAME
DB_USER
DB_PASSWORD
```

Для локального запуска без Docker можно использовать файл `.env`.

---

## Требования

- Docker
- Docker Compose