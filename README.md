# 🏗️ PingTower

**Современная система мониторинга сайтов на микросервисной архитектуре**

Высокопроизводительное решение для отслеживания доступности веб-ресурсов с real-time уведомлениями и полной API документацией.

## ✨ Возможности

- 🔍 **Мониторинг в реального времени** - проверка доступности сайтов
- 📧 **Smart уведомления** - email алерты через MailerSend SMTP
- 🔐 **JWT авторизация** - безопасная аутентификация пользователей  
- 📊 **Аналитика** - ClickHouse для хранения логов и метрик
- 🌐 **REST API** - полная документация в Swagger UI
- 🐳 **Docker-ready** - простое развертывание в контейнерах
- ⚡ **Kafka очереди** - асинхронная обработка уведомлений
- 🎯 **Микросервисы** - масштабируемая архитектура

## 📋 Содержание

- [Архитектура](#архитектура)
- [Быстрый старт](#-быстрый-старт)  
- [Использование](#использование)
- [API Документация](#-api-документация)
- [Настройка](#️-настройка)
- [Разработка](#разработка)
- [Отладка](#-отладка-и-решение-проблем)
- [Технологии](#️-технологии)

## Архитектура

Система построена на микросервисной архитектуре и состоит из 5 основных сервисов:

![Архитектура системы](architecture.png)

### 🌐 Микросервисы

- **🌐 API Service** (8080) - Gateway, точка входа для клиентов, JWT авторизация, запуск Cron для систематического пинга сайтов пользователей.
- **🔐 Auth Service** (8081) - генерация, валидация и обновление JWT токенов
- **📡 Ping Service** (8082) - проверка доступности сайтов с таймаутами
- **🗄️ DB Service** (8083) - управление данными, PostgreSQL и ClickHouse
- **📧 Notification Service** (8084) - email уведомления через Kafka

### 🗄️ Базы данных

- **PostgreSQL** (5432) - пользователи, сайты, конфигурация
- **ClickHouse** (8123/9000) - логи проверок, аналитика
- **Redis** (6379) - refresh токены, кэширование

### 📨 Очередь сообщений

- **Apache Kafka** (19092) - асинхронные уведомления
- **Zookeeper** (2181) - координация Kafka

### 🔧 Мониторинг и UI

- **Kafka UI** (8085) - мониторинг очередей сообщений  
- **Swagger UI** (8080/swagger) - интерактивная документация API

## 🚀 Быстрый старт

### Требования
- **Docker** (20.10+)
- **Docker Compose** (2.0+)
- **8 ГБ RAM** (рекомендуется)

### 1️⃣ Клонирование и запуск
```bash
# Скачать проект
git clone <repository-url>
cd PingTower

# Запустить все сервисы
docker-compose up -d --build

# Дождаться загрузки (1-2 минуты)
docker-compose logs -f api_service
```

### 2️⃣ Проверка работоспособности
```bash
# Health check основных сервисов
curl http://localhost:8080/health          # API Gateway
curl http://localhost:8081/health          # Auth Service  
curl http://localhost:8084/health          # Notification Service

# Доступность интерфейсов
open http://localhost:8080/swagger         # Swagger UI
open http://localhost:8085                 # Kafka UI
```

### 3️⃣ Готово! 
Система запущена и готова к использованию через веб-интерфейс и API.

## 🌐 Веб-интерфейс

### Доступные интерфейсы:
- **Swagger UI** - http://localhost:8080/swagger - интерактивная документация и тестирование API
- **Kafka UI** - http://localhost:8085 - мониторинг очередей сообщений

## 📖 API Документация

### 🌟 Swagger UI

**Полная интерактивная документация доступна по адресу:**
```
http://localhost:8080/swagger/
```

Swagger UI содержит:
- ✅ **22 эндпоинта** всех 5 микросервисов
- ✅ **Интерактивное тестирование** API
- ✅ **Примеры запросов** и ответов
- ✅ **Схемы данных** и валидация
- ✅ **Авторизация** через Bearer токены

### 🔑 Основные эндпоинты

#### API Service (8080)
- `POST /register` - регистрация пользователя
- `POST /login` - авторизация
- `GET /checkers` - список сайтов пользователя
- `POST /checkers` - добавить сайт для мониторинга
- `GET /checker/{id}` - логи конкретного сайта
- `POST /pingAll` - запустить проверку всех сайтов
- `GET /health` - health check

#### Auth Service (8081)
- `POST /generate` - генерация JWT токенов
- `POST /validate` - валидация токена
- `POST /refresh` - обновление токенов

#### DB Service (8083)
- `GET /all-users-sites` - все пользователи и сайты
- `POST /ping` - сохранение лога пинга
- `GET /user/{id}/email` - получить email пользователя

### 🔐 Авторизация

Все защищенные эндпоинты требуют JWT токен:
```bash
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## ⚙️ Настройка

### 🔧 Переменные окружения

Система использует переменные окружения в `docker-compose.yml` и `.env` файлах:

#### Auth Service
```env
SERVER_PORT=8081
REDIS_URL=redis:6379
JWT_SECRET=your-super-secret-jwt-key-here
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d
```

#### Notification Service
```env
MAILERSEND_API_KEY=mlsn.56553fafd04e982332a4ef9003ef6fc0345de9fcea0473122e482322f2b83938
FROM_EMAIL=MS_mIKkkJ@test-xkjn41m6jz04z781.mlsender.net
FROM_NAME=PingTower Alert System
SMTP_HOST=smtp.mailersend.net
SMTP_PORT=587
SMTP_USERNAME=MS_mIKkkJ@test-xkjn41m6jz04z781.mlsender.net
SMTP_PASSWORD=mssp.SFuGnXj.zr6ke4n0z39lon12.CtzsE9m
KAFKA_BROKERS=kafka1:29092
KAFKA_TOPIC=notification-alerts
KAFKA_CONSUMER_GROUP=notification-service
HEALTH_PORT=8084
```

#### База данных
```env
# PostgreSQL
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=ping_db

# ClickHouse
CLICKHOUSE_USER=default
CLICKHOUSE_PASSWORD=
CLICKHOUSE_DB=ping_db
```

### 📧 Настройка SMTP

Для получения email уведомлений настройте MailerSend:

1. Зарегистрируйтесь на [MailerSend](https://mailersend.com)
2. Получите API ключ и SMTP данные
3. Обновите переменные в `docker-compose.yml`:
   - `MAILERSEND_API_KEY`
   - `FROM_EMAIL` 
   - `SMTP_USERNAME`
   - `SMTP_PASSWORD`

### 🔒 Безопасность

**⚠️ Важно:** Перед продакшеном обновите:
- `JWT_SECRET` - уникальный секретный ключ
- Пароли баз данных
- SMTP ключи

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

## 🔧 Отладка и решение проблем

### 📋 Диагностика системы
```bash
# Статус всех контейнеров
docker-compose ps

# Использование ресурсов
docker stats --no-stream

# Логи конкретного сервиса
docker-compose logs -f [api_service|auth_service|ping_service|db_service|notification_service]

# Последние 50 строк логов всех сервисов
docker-compose logs --tail=50
```

### 🗄️ Проверка баз данных
```bash
# PostgreSQL
docker exec pingtower-postgres_db-1 psql -U postgres -d ping_db -c "SELECT version();"
docker exec pingtower-postgres_db-1 psql -U postgres -d ping_db -c "\dt"

# ClickHouse
docker exec pingtower-clickhouse_db-1 clickhouse-client --query "SELECT version()"
docker exec pingtower-clickhouse_db-1 clickhouse-client --query "SHOW TABLES"

# Redis
docker exec pingtower-redis-1 redis-cli ping
docker exec pingtower-redis-1 redis-cli info stats
```

### 📨 Проверка Kafka
```bash
# Список топиков
docker exec pingtower-kafka1-1 kafka-topics --bootstrap-server localhost:29092 --list

# Проверка сообщений в топике
docker exec pingtower-kafka1-1 kafka-console-consumer \
  --bootstrap-server localhost:29092 \
  --topic notification-alerts \
  --from-beginning --max-messages 10

# Kafka UI: http://localhost:8085
```

### 🚨 Частые проблемы

#### 1. Сервисы не запускаются
```bash
# Очистить и пересобрать
docker-compose down -v
docker-compose build --no-cache
docker-compose up -d

# Проверить ресурсы
docker system df
docker system prune -f
```

#### 2. База данных недоступна
```bash
# Перезапустить PostgreSQL
docker-compose restart postgres_db

# Проверить подключение
docker-compose exec postgres_db pg_isready -U postgres
```

#### 3. Email не отправляются  
```bash
# Проверить логи notification service
docker-compose logs notification_service

# Проверить конфигурацию SMTP в docker-compose.yml
grep -A10 MAILERSEND docker-compose.yml
```

#### 4. JWT токены не работают
```bash
# Перезапустить auth service
docker-compose restart auth_service

# Проверить Redis
docker exec pingtower-redis-1 redis-cli keys "*"
```

### 🔄 Перезапуск системы
```bash
# Мягкий перезапуск
docker-compose restart

# Полный перезапуск с очисткой
docker-compose down
docker-compose up -d

# Перезапуск с пересборкой
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## 🛠️ Технологии

### Backend
- **Go 1.24** - основной язык программирования
- **Microservices** - архитектурный паттерн
- **JWT** - аутентификация и авторизация
- **REST API** - HTTP интерфейс

### Базы данных
- **PostgreSQL 15** - пользователи, сайты, конфигурация
- **ClickHouse 23.8** - логи проверок, временные ряды
- **Redis 7** - refresh токены, кэширование

### Очередь сообщений
- **Apache Kafka 7.4.0** - асинхронная обработка уведомлений
- **Zookeeper** - координация Kafka кластера

### Инфраструктура
- **Docker & Docker Compose** - контейнеризация
- **Swagger/OpenAPI 3.0.3** - документация API
- **SMTP (MailerSend)** - доставка email

### Мониторинг и отладка
- **Kafka UI** - http://localhost:8085 - мониторинг очередей
- **Swagger UI** - http://localhost:8080/swagger - тестирование API  
- **Health checks** - статус всех сервисов
- **Структурированное логирование** - отладка и аналитика

### 📊 Доступные интерфейсы

| Сервис | URL | Описание |
|--------|-----|----------|
| **API Gateway** | http://localhost:8080 | Основной API |
| **Swagger UI** | http://localhost:8080/swagger | Документация API |
| **Kafka UI** | http://localhost:8085 | Мониторинг очередей |
| **PostgreSQL** | localhost:5432 | База данных |
| **ClickHouse** | localhost:8123 | Аналитическая БД |
| **Redis** | localhost:6379 | Кэш |
