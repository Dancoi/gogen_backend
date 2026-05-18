# Развертывание Gogen Backend

## Таблица содержания
1. [Требования](#требования)
2. [Docker Compose (Рекомендуется)](#docker-compose-рекомендуется)
3. [Локальная разработка](#локальная-разработка)
4. [Production развертывание](#production-развертывание)

---

## Требования

### Для Docker Compose
- Docker >= 20.10
- Docker Compose >= 1.29
- (опционально) Make

### Для локальной разработки
- Go >= 1.25
- PostgreSQL >= 15
- Goose (для миграций)
- SQLC (для генерации кода)

---

## Docker Compose (Рекомендуется)

### Быстрый старт (3 команды)

```bash
# 1. Клонируйте репозиторий
git clone <repository-url>
cd gogen_backend

# 2. Запустите контейнеры
docker-compose up -d

# 3. Проверьте статус
docker-compose ps
```

### Проверка

```bash
# Проверьте API
curl http://localhost:8080/health

# Ожидается ответ:
# {"status":"ok"}
```

### Просмотр логов

```bash
# Все сервисы
docker-compose logs -f

# Только backend
docker-compose logs -f gogen_backend

# Только БД
docker-compose logs -f postgres
```

### Остановка

```bash
# Остановить контейнеры
docker-compose stop

# Остановить и удалить
docker-compose down

# Удалить также volumes (ОСТОРОЖНО - удалит БД!)
docker-compose down -v
```

### Переменные окружения

Отредактируйте `.env` файл:

```bash
# Database (внутри контейнера используйте "postgres" как host)
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gogen

# Server
SERVER_PORT=8080

# JWT (сгенерируйте сильный ключ в production!)
JWT_SECRET=your-secret-key-change-in-production
```

---

## Локальная разработка

### Установка зависимостей

```bash
# Go modules
go mod download

# SQLC (для генерации SQL кода)
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Goose (для миграций)
go install github.com/pressly/goose/v3/cmd/goose@latest

# Gin веб-фреймворк и другие зависимости
go mod tidy
```

### Запуск PostgreSQL (вариант 1 - Docker)

```bash
docker run --name gogen_postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=gogen \
  -p 5432:5432 \
  -v postgres_data:/var/lib/postgresql/data \
  postgres:18-alpine
```

### Запуск PostgreSQL (вариант 2 - локально)

Установите PostgreSQL и создайте БД:

```bash
createdb -U postgres gogen
```

### Миграции

```bash
# Установите переменные окружения
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="host=localhost port=5432 user=postgres password=postgres dbname=gogen sslmode=disable"
export GOOSE_MIGRATION_DIR=./sql/migrations

# Запустите миграции
goose up

# Проверьте статус
goose status
```

### Генерация SQLC кода

```bash
# Убедитесь что БД с миграциями готова
sqlc generate
```

### Запуск backend

```bash
# Создайте .env файл
cp .env.example .env

# Запустите
go run ./cmd/main/main.go

# Или соберите бинарник
go build -o server ./cmd/main/main.go
./server
```

### Тестирование API

```bash
# Health check
curl http://localhost:8080/health

# Регистрация
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "john_doe",
    "password": "SecurePassword123"
  }'
```

---

## Production развертывание

### Рекомендуемая архитектура

```
┌─────────────┐
│   Nginx     │ (reverse proxy, SSL/TLS)
└──────┬──────┘
       │
┌──────▼──────────────────────┐
│   Docker Compose            │
├─────────────────────────────┤
│  ┌────────────┐              │
│  │ Go Backend │              │
│  │ (port 8080)│              │
│  └────────────┘              │
│                              │
│  ┌────────────┐              │
│  │ PostgreSQL │              │
│  │ (port 5432)│              │
│  └────────────┘              │
└─────────────────────────────┘
```

### 1. Подготовка сервера

```bash
# Обновите систему
sudo apt update && sudo apt upgrade -y

# Установите Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Добавьте пользователя в docker группу
sudo usermod -aG docker $USER

# Установите Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### 2. Переменные окружения

Создайте `.env` файл с production переменными:

```bash
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=$(openssl rand -hex 16)  # Сгенерируйте сильный пароль
DB_NAME=gogen

# Server
SERVER_PORT=8080

# JWT (ОБЯЗАТЕЛЬНО сгенерируйте!)
JWT_SECRET=$(openssl rand -hex 32)

# Environment
GO_ENV=production
```

### 3. docker-compose.yml для production

Отредактируйте файл с добавлением лимитов ресурсов:

```yaml
services:
  gogen_backend:
    # ... остальная конфигурация ...
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
    # ... остальная конфигурация ...
```

### 4. Nginx конфигурация

Создайте `/etc/nginx/sites-available/gogen`:

```nginx
upstream backend {
    server 127.0.0.1:8080;
}

server {
    listen 80;
    server_name api.example.com;
    
    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.example.com;
    
    ssl_certificate /etc/letsencrypt/live/api.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.example.com/privkey.pem;
    
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    
    client_max_body_size 10M;
    
    location / {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
    }
}
```

Активируйте:

```bash
sudo ln -s /etc/nginx/sites-available/gogen /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 5. SSL сертификат (Let's Encrypt)

```bash
# Установите Certbot
sudo apt install certbot python3-certbot-nginx -y

# Получите сертификат
sudo certbot certonly --nginx -d api.example.com
```

### 6. Запуск на production

```bash
# Создайте папку для приложения
mkdir -p /opt/gogen
cd /opt/gogen

# Скопируйте файлы
git clone <repository> .
cp .env.example .env

# Отредактируйте .env с production переменными
nano .env

# Запустите контейнеры в фоне
docker-compose up -d

# Проверьте статус
docker-compose ps
docker-compose logs
```

### 7. Мониторинг и backups

```bash
# Создайте скрипт backup (backup.sh)
#!/bin/bash
BACKUP_DIR="/backups/gogen"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR
docker-compose exec -T postgres pg_dump -U postgres gogen | gzip > $BACKUP_DIR/gogen_$DATE.sql.gz

# Удалите старые backups (старше 30 дней)
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete

# Добавьте в crontab
# 0 2 * * * /opt/gogen/backup.sh
```

### 8. Обновление приложения

```bash
# Перейдите в директорию
cd /opt/gogen

# Обновите код
git pull

# Пересоберите контейнер
docker-compose build --no-cache gogen_backend

# Перезапустите
docker-compose up -d

# Проверьте логи
docker-compose logs -f gogen_backend
```

---

## Troubleshooting

### Контейнер не запускается

```bash
# Проверьте логи
docker-compose logs gogen_backend

# Проверьте конфигурацию
docker-compose config

# Пересоберите без кэша
docker-compose build --no-cache
```

### Проблема подключения к БД

```bash
# Проверьте что postgres запущен
docker-compose ps

# Проверьте переменные окружения
docker-compose exec gogen_backend env | grep DB

# Проверьте network
docker network ls
docker network inspect gogen_network
```

### Сбросить БД

```bash
# ВНИМАНИЕ: Это удалит все данные!
docker-compose down -v
docker-compose up -d
```

---

## Дополнительные ресурсы

- [Docker документация](https://docs.docker.com/)
- [Docker Compose документация](https://docs.docker.com/compose/)
- [Go на Docker](https://docs.docker.com/language/golang/)
- [PostgreSQL Docker](https://hub.docker.com/_/postgres)
