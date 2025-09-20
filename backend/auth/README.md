# Auth Service

Микросервис аутентификации и авторизации для PingTower.

## API Endpoints

### POST /register
Регистрация нового пользователя
```bash
curl -X POST http://localhost:8081/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

### POST /login
Вход в систему
```bash
curl -X POST http://localhost:8081/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

### POST /validate
Проверка токена (для других сервисов)
```bash
# Через заголовок Authorization
curl -X POST http://localhost:8081/validate \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Или через JSON body
curl -X POST http://localhost:8081/validate \
  -H "Content-Type: application/json" \
  -d '{"token": "YOUR_JWT_TOKEN"}'
```

## Переменные окружения
Скопируйте `.env.example` в `.env` и настройте:
- `JWT_SECRET` - секретный ключ для JWT
- `JWT_EXPIRY_HOURS` - время жизни токена в часах
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - настройки PostgreSQL
- `SERVER_PORT` - порт сервера

## Запуск
```bash
# Локально
go run cmd/main.go

# Docker
docker build -t auth-service .
docker run -p 8081:8081 auth-service
```