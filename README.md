# 📝 Notes API

REST API для управления заметками с JWT аутентификацией, написанный на Go.

## 🛠 Технологии

- **Go** — основной язык
- **PostgreSQL** — база данных
- **Gorilla Mux** — роутер
- **JWT** — аутентификация
- **bcrypt** — хэширование паролей

## 📁 Структура проекта

```
notes-api/
├── main.go
├── config/
│   └── database.go
├── models/
│   ├── user.go
│   └── note.go
├── handlers/
│   ├── auth.go
│   └── notes.go
├── middleware/
│   └── auth.go
└── routes/
    └── routes.go
```

## ⚙️ Установка и запуск

### 1. Клонируй репозиторий

```bash
git clone https://github.com/твой-username/notes-api.git
cd notes-api
```

### 2. Создай базу данных PostgreSQL

```sql
CREATE DATABASE notesdb;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT,
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 3. Настрой подключение к БД

В файле `config/database.go` укажи свои данные:

```go
connStr := "host=localhost port=5432 user=postgres password=YOUR_PASSWORD dbname=notesdb sslmode=disable"
```

### 4. Установи зависимости и запусти

```bash
go mod tidy
go run main.go
```

Сервер запустится на `http://localhost:8080`

## 🔑 Аутентификация

API использует **JWT токены**. После логина добавляй токен в заголовок каждого запроса:

```
Authorization: Bearer ваш_токен
```

## 📡 Эндпоинты

### Публичные (без токена)

| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/ping` | Проверка работы сервера |
| POST | `/register` | Регистрация |
| POST | `/login` | Вход, получение токена |

### Защищённые (нужен токен)

| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/notes` | Получить все заметки |
| GET | `/notes/{id}` | Получить заметку по ID |
| POST | `/notes` | Создать заметку |
| PUT | `/notes/{id}` | Обновить заметку |
| DELETE | `/notes/{id}` | Удалить заметку |

## 📋 Примеры запросов

### Регистрация

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secret123"}'
```

### Логин

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secret123"}'
```

Ответ:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Создать заметку

```bash
curl -X POST http://localhost:8080/notes \
  -H "Authorization: Bearer ВАШ_ТОКЕН" \
  -H "Content-Type: application/json" \
  -d '{"title":"Заголовок","content":"Текст заметки"}'
```

### Получить все заметки

```bash
curl http://localhost:8080/notes \
  -H "Authorization: Bearer ВАШ_ТОКЕН"
```

### Обновить заметку

```bash
curl -X PUT http://localhost:8080/notes/1 \
  -H "Authorization: Bearer ВАШ_ТОКЕН" \
  -H "Content-Type: application/json" \
  -d '{"title":"Новый заголовок","content":"Новый текст"}'
```

### Удалить заметку

```bash
curl -X DELETE http://localhost:8080/notes/1 \
  -H "Authorization: Bearer ВАШ_ТОКЕН"
```

## 📊 Коды ответов

| Код | Значение |
|-----|----------|
| 200 | Успешно |
| 201 | Создано |
| 204 | Удалено |
| 400 | Неверный запрос |
| 401 | Не авторизован |
| 404 | Не найдено |
| 409 | Username занят |
| 500 | Ошибка сервера |
