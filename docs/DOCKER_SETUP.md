# Docker Setup Guide

## Предварительные требования

- Docker ([установка](https://docs.docker.com/get-docker/))
- Docker Compose ([установка](https://docs.docker.com/compose/install/))

## Быстрый старт

### 1. Подготовка переменных окружения

```bash
# Скопируйте .env.example в .env
cp .env.example .env

# (опционально) Отредактируйте переменные окружения
# nano .env
```

### 2. Запуск контейнеров

```bash
# Запуск всех сервисов в режиме foreground
docker-compose up

# Запуск в фоновом режиме (daemon)
docker-compose up -d

# С пересборкой образов
docker-compose up --build
```

### 3. Проверка статуса

```bash
# Просмотр запущенных контейнеров
docker-compose ps

# Просмотр логов
docker-compose logs -f

# Логи конкретного сервиса
docker-compose logs -f gogen_backend
```

### 4. Остановка сервисов

```bash
# Остановить все контейнеры
docker-compose stop

# Остановить и удалить контейнеры
docker-compose down

# Остановить и удалить всё (включая volumes)
docker-compose down -v
```

## Структура Docker Compose

### Сервисы

#### PostgreSQL (postgres)
- **Image**: postgres:18-alpine
- **Port**: 5432 (по умолчанию)
- **Volume**: postgres_data (persistent storage)
- **Health Check**: pg_isready

#### Go Backend (gogen_backend)
- **Build**: Dockerfile (multi-stage)
- **Зависит от**: postgres (healthcheck)
- **Port**: 8080 (по умолчанию)
- **Restart Policy**: unless-stopped

## Сборка вручную

```bash
# Сборка только образа backend
docker-compose build gogen_backend

# Сборка с использованием кэша
docker-compose build --no-cache gogen_backend
```

## Миграции БД

Миграции запускаются автоматически при первом старте контейнера (благодаря healthcheck).

Если нужно запустить вручную:

```bash
# Внутри контейнера
docker-compose exec gogen_backend goose -dir ./sql/migrations postgres "postgres://postgres:postgres@postgres:5432/gogen?sslmode=disable" up
```

## Доступ к сервисам

### API Backend
```bash
curl http://localhost:8080/health
```

### PostgreSQL
```bash
# Подключение с помощью psql
docker-compose exec postgres psql -U postgres -d gogen

# Или через pgAdmin (если добавите)
```

## Проблемы и решения

### Ошибка "port is already in use"

```bash
# Найти процесс использующий порт (Linux/Mac)
lsof -i :8080

# Найти процесс использующий порт (Windows)
netstat -ano | findstr :8080

# Изменить порты в .env
# DB_PORT=5433
# SERVER_PORT=8081
```

### Контейнер падает при старте

```bash
# Посмотрите логи
docker-compose logs gogen_backend

# Проверьте переменные окружения в .env
# Убедитесь что база готова (health check)
```

### Очистка старых данных

```bash
# Удалить все volume'ы
docker-compose down -v

# Заново запустить
docker-compose up
```

## Production рекомендации

1. **Используйте сильный JWT_SECRET**:
   ```bash
   openssl rand -hex 32
   ```

2. **Настройте HTTPS** (используйте reverse proxy like Nginx)

3. **Используйте отдельные .env файлы** для разных окружений:
   - `.env.production`
   - `.env.staging`

4. **Установите лимиты ресурсов** в docker-compose.yml:
   ```yaml
   services:
     gogen_backend:
       deploy:
         resources:
           limits:
             cpus: '1'
             memory: 512M
           reservations:
             cpus: '0.5'
             memory: 256M
   ```

5. **Мониторинг и логирование**:
   - Используйте docker logs driver (json-file, syslog, etc.)
   - Интегрируйте с системой мониторинга

## Dockerfile Stages

### Stage 1: Builder
- Использует golang:1.25-alpine
- Компилирует приложение в статический бинарник
- Большой размер, но используется только для сборки

### Stage 2: Runtime
- Использует alpine:latest
- Копирует только скомпилированный бинарник
- Минимальный размер финального образа (~100MB)

## Команды для разработки

```bash
# Пересобрать только backend
docker-compose build gogen_backend

# Запустить и посмотреть логи
docker-compose up gogen_backend --build

# Выполнить команду внутри контейнера
docker-compose exec gogen_backend go test ./...

# Интерактивный shell в контейнере
docker-compose exec gogen_backend sh
```

## Сетевое взаимодействие

Все сервисы в одной сети (gogen_network), поэтому:
- Backend может подключиться к базе по hostname: `postgres`
- Внешние клиенты подключаются по `localhost` на экспортированные порты
