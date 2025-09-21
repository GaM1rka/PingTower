# PingTower

**Система мониторинга сайтов**

Простое решение для отслеживания доступности сайтов с уведомлениями на email.

## Возможности

- Мониторинг сайтов в реальном времени
- Email уведомления при падении
- Авторизация пользователей
- Веб-интерфейс для управления

## Архитектура

Система состоит из нескольких сервисов:

- **API** (8080) - главный сервис
- **Auth** (8081) - авторизация 
- **Ping** (8082) - проверка сайтов
- **Notifications** (8084) - отправка email
- **PostgreSQL** (5432) - база данных
- **Redis** (6379) - кеш и сессии

## Запуск

Требования: Docker

```bash
# Скачать проект
git clone <repository-url>
cd PingTower

# Запустить все сервисы
docker-compose up -d

# Проверить что работает
curl http://localhost:8080/health
```

## Использование

### Регистрация пользователя
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'
```

### Добавление сайта для мониторинга
```bash
curl -X POST http://localhost:8080/checkers \
  -H "Authorization: Bearer ВАШ_ТОКЕН" \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com","name":"Мой сайт"}'
```

### Запуск проверки всех сайтов
```bash
curl -X POST http://localhost:8080/pingAll
```

## API

### Основные эндпоинты

- `POST /register` - регистрация
- `POST /login` - вход  
- `GET /checkers` - список сайтов
- `POST /checkers` - добавить сайт
- `POST /pingAll` - проверить все сайты

### Авторизация

Все запросы кроме регистрации и входа требуют JWT токен в заголовке:
```
Authorization: Bearer ВАШ_ТОКЕН
```

## Настройка

Основные переменные окружения в docker-compose.yml:

- `MAILERSEND_API_KEY` - ключ для отправки email
- `SMTP_USERNAME` - логин SMTP
- `SMTP_PASSWORD` - пароль SMTP  
- `JWT_SECRET` - секрет для JWT токенов

## Разработка

Запуск отдельных сервисов:

```bash
# API сервис
cd backend/api_service && go run cmd/main.go

# Авторизация  
cd backend/auth && go run cmd/main.go

# Проверка сайтов
cd backend/ping_service && go run cmd/main.go

# Уведомления
cd backend/notification_service && go run cmd/main.go
```

## Проблемы

Если что-то не работает:

```bash
# Посмотреть логи
docker-compose logs [название-сервиса]

# Перезапустить сервис
docker-compose restart [название-сервиса]

# Проверить базы данных
docker exec pingtower-postgres_db-1 psql -U postgres -d ping_db -c "SELECT version();"
docker exec pingtower-redis-1 redis-cli ping
```

## Технологии

- Go - основной язык
- PostgreSQL - основная БД
- Redis - кеш и сессии
- Kafka - очередь сообщений
- Docker - контейнеризация
- MailerSend - отправка email