# 📰 News API

Простое и аккуратное REST API для управления новостями с аутентификацией по JWT. Готово для запуска в Docker и локальной
разработки.

-----

## ✨ Возможности

* Регистрация и вход пользователей
* Роли и доступ: `admin`, `editor`, `user`
* CRUD для новостей (создание/чтение/обновление/удаление)
* Пагинация, поиск
* Хранение сессий/рефреш‑токенов в Redis
* Миграции БД, готовые curl‑примеры

-----

## 🧰 Стек

* **Go** (Gorilla/mux)
* **PostgreSQL**
* **Redis**
* **Docker + Docker Compose**
* **JWT‑аутентификация** (Access + Refresh)

-----

## 📁 Структура проекта (пример)

```
.
├── Dockerfile
├── Makefile
├── README.md
├── cmd
│   └── main.go
├── docker-compose.yml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│   ├── app
│   ├── config
│   ├── database
│   ├── dto
│   ├── http
│   ├── middleware
│   ├── models
│   ├── repository
│   └── service
├── migrations
│   ├── 20250823164042_create_users_table.sql
│   └── 20250823164047_create_news_table.sql
├── pkg
│   ├── logger
│   ├── password
│   ├── redis
│   └── token
└── utils
    └── writer.go
```

-----

## 🚀 Быстрый старт (Docker)

### 1\) Клонируйте репозиторий

```bash
git clone https://github.com/your-username/news-api.git
cd news-api
```

### 2\) Соберите и запустите сервисы

```bash
docker-compose up --build -d
```

-----

-----

## 🔐 Аутентификация и сессии

* **Access Token** (короткоживущий) — передаётся в `Authorization: Bearer <token>`
* **Refresh Token** (долго живёт) — хранится в Redis.
* **Logout** — инвалидируем связанный refresh (и при необходимости помещаем access в blacklist до истечения TTL).

-----

## 📚 REST API (основные ручки)

Базовый префикс: `/api`

### Аутентификация

* `POST /api/register` — регистрация пользователя (`email`, `password`)
* `POST /api/login` — вход, возвращает пары токенов `{access, refresh}`
* `POST /api/logout` — выход

### Новости

* `GET    /api/news` — список с пагинацией/поиском
* `GET    /api/news/{id}` — получить новость
* `POST   /api/news` — создать (роль: `editor`/`admin`)
* `PUT    /api/news/{id}` — обновить (роль: `editor`/`admin` и/или автор)
* `DELETE /api/news/{id}` — удалить (роль: `admin`)

-----

### ⚙️ Конфигурация

Создайте файл `.env` в корневой директории проекта, используя следующий пример:

```env
# Server Configuration
SERVER_PORT=8080

# Database Configuration (PostgreSQL)
DB_HOST=news-api-db
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=news_db
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=supersecret
JWT_EXPIRATION_HOURS=1

# Redis Configuration
REDIS_ADDR=news-api-redis:6379
REDIS_USERNAME=
REDIS_PASSWORD=

# Logger
LOG_LEVEL=info
```

-----

## Можете начинать работать с помощью swagger

* `http://localhost:8080/swagger/index.html#/`

-----