# Setup Guide

## Запуск Docker Compose

```bash
# Создать .env файл из шаблона
cp .env.example .env

# Поднять контейнеры
docker-compose up -d

# Проверить статус
docker-compose ps

# Остановить контейнеры
docker-compose down
```

## Доступ к БД

- **PostgreSQL**: localhost:5432
- **PgAdmin**: http://localhost:5050
  - Email: admin@admin.com
  - Password: admin

## Инструменты для миграций

### 1. **Migrate (Рекомендуется)** ⭐
Самый популярный инструмент для миграций в Go.

**Установка:**
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

**Использование:**
```bash
# Создать миграцию
migrate create -ext sql -dir migrations -seq create_users_table

# Применить миграции
migrate -path migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up

# Откатить на 1 версию
migrate -path migrations -database "postgresql://..." down 1

# Получить текущую версию
migrate -path migrations -database "postgresql://..." version
```

**Пример структуры миграций:**
```
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_posts_table.up.sql
└── 000002_create_posts_table.down.sql
```

### 2. **Goose**
Альтернатива Migrate, поддерживает Go и SQL миграции.

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest

# Создать миграцию
goose create create_users_table sql

# Применить миграции
goose postgres "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME" up
```

### 3. **SQLC** (для type-safe SQL)
Генерирует Go код из SQL запросов.

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

## Рекомендация

Использую **Migrate** - он:
- ✅ Простой и надежный
- ✅ Широко используется в production
- ✅ Поддерживает откаты
- ✅ Версионирование миграций
- ✅ Хорошая документация

Если нужен type-safe SQL - комбинируй **Migrate** для структуры БД и **SQLC** для запросов.
