# Инструкции по запуску SQLC и Goose

## 1. Установка инструментов

```bash
# Goose (для миграций)
go install github.com/pressly/goose/v3/cmd/goose@latest

# SQLC (для генерации Go кода)
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# pgx драйвер (для работы с PostgreSQL)
go get github.com/jackc/pgx/v5
```

## 2. Запуск Docker Compose

```bash
# Поднять PostgreSQL и PgAdmin
docker-compose up -d

# Проверить статус
docker-compose ps
```

## 3. Запуск миграций Goose

```bash
# Запустить все миграции
goose -dir sql/migrations postgres "postgres://postgres:postgres@localhost:5432/gogen?sslmode=disable" up

# Проверить текущую версию
goose -dir sql/migrations postgres "postgres://postgres:postgres@localhost:5432/gogen?sslmode=disable" version

# Откатить последнюю миграцию
goose -dir sql/migrations postgres "postgres://postgres:postgres@localhost:5432/gogen?sslmode=disable" down

# Откатить все миграции
goose -dir sql/migrations postgres "postgres://postgres:postgres@localhost:5432/gogen?sslmode=disable" down-to 0
```

## 4. Генерация SQLC кода

```bash
# После создания миграций и SQL запросов
sqlc generate

# Это создаст файлы в internal/sqlc/
# - models.go (структуры данных)
# - queries.go (функции для запросов)
# - db.go (интерфейс Querier)
```

## 5. Использование в коде

```go
import (
    "context"
    "github.com/Dancoi/gogen_backend/internal/sqlc"
    "github.com/Dancoi/gogen_backend/pkg/db"
    "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
    ctx := context.Background()
    
    // Подключаемся к БД
    pool, err := db.New(ctx, "postgres://postgres:postgres@localhost:5432/gogen?sslmode=disable")
    if err != nil {
        panic(err)
    }
    defer pool.Close()
    
    // Создаём Queries объект
    queries := sqlc.New(pool)
    
    // Используем запросы
    user, err := queries.GetUserByEmail(ctx, "user@example.com")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("User: %v\n", user)
}
```

## 6. Структура файлов

```
sql/migrations/
├── 000001_init_schema.sql       # Создание таблиц (Up и Down в одном)
└── 000002_add_indexes.sql       # Добавление индексов (Up и Down в одном)

sql/queries/
├── users.sql                    # Запросы для users
├── subscriptions.sql            # Запросы для subscriptions
├── api_tokens.sql              # Запросы для api_tokens
├── user_sessions.sql           # Запросы для user_sessions
├── audit_logs.sql              # Запросы для audit_logs
├── token_usage_logs.sql        # Запросы для token_usage_logs
├── roles.sql                   # Запросы для roles
└── subscription_plans.sql      # Запросы для subscription_plans

internal/sqlc/                  # ГЕНЕРИРУЕТСЯ АВТОМАТИЧЕСКИ
├── models.go                   # Структуры (User, Subscription, etc)
├── db.go                       # Интерфейс Querier и коннекшен
├── users.sql.go               # Сгенерированные функции
├── subscriptions.sql.go
├── api_tokens.sql.go
└── ...
```

## 7. Полезные команды

```bash
# Просмотр всех таблиц в БД
psql -h localhost -U postgres -d gogen -c "\dt"

# Просмотр структуры таблицы
psql -h localhost -U postgres -d gogen -c "\d users"

# Проверить валидность SQLC конфига
sqlc compile

# Генерировать с verbose выводом
sqlc generate --verbose
```

## 8. Переменные окружения

Создать `.env` файл:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gogen
```

И использовать в коде:
```go
cfg := config.LoadConfig()
dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
    cfg.DBUser,
    cfg.DBPassword,
    cfg.DBHost,
    cfg.DBPort,
    cfg.DBName,
)
```
