# GoGen API Documentation 📚

Полная документация по API личного кабинета для получения токена для консольного инструмента на Go.

## 📖 Содержание

### 1. **[API.md](./API.md)** - Полная документация API
Описание всех endpoints, параметров запросов и ответов:
- ✅ Health Check
- 🔐 Регистрация и авторизация
- 🎫 Управление API токенами
- 📋 Структура ответов и ошибок
- 🔑 Аутентификация (JWT и API токены)

**Используй этот файл для:**
- Интеграции с фронтенд приложением
- Понимания всех доступных endpoints
- Изучения формата запросов и ответов

---

### 2. **[EXAMPLES.md](./EXAMPLES.md)** - Примеры использования
Готовые примеры кода на разных языках и фреймворках:
- 🟦 TypeScript / React
- 🟨 JavaScript / Vanilla
- 🐍 Python
- 🔷 Go
- 💻 Bash / cURL
- 📮 Postman Collection

**Используй этот файл для:**
- Быстрого старта интеграции
- Copy-paste готовых решений
- Понимания как использовать API в твоем проекте

---

### 3. **[DATA_MODELS.md](./DATA_MODELS.md)** - Структуры данных
Описание всех моделей данных и структур БД:
- 👤 User - Пользователь
- 🎫 ApiToken - API Токен
- 📅 Subscription - Подписка
- 📋 AuditLog - Логи аудита
- 🔐 Role - Роли пользователя
- И много других...

**Используй этот файл для:**
- Понимания структуры данных
- Проектирования фронтенд компонентов
- Работы с JSON ответами API

---

### 4. **[FAQ.md](./FAQ.md)** - Часто задаваемые вопросы
Ответы на типичные вопросы разработчиков:
- 🤔 Как сохранить JWT токен?
- 🔑 Как использовать API токены?
- 🆘 Как решить ошибки?
- 🛡️ Вопросы безопасности
- 🧪 Тестирование и разработка

**Используй этот файл для:**
- Быстрого поиска ответа на свой вопрос
- Решения проблем
- Понимания best practices

---

## 🚀 Быстрый старт

### 1. Установка
```bash
# Клонировать репозиторий
git clone https://github.com/Dancoi/gogen_backend.git
cd gogen_backend

# Установить зависимости
go mod download

# Запустить сервер
go run ./cmd/main/main.go
```

### 2. Локальное тестирование
```bash
# Сервер запустится на http://localhost:8080

# Проверить доступность
curl http://localhost:8080/health
```

### 3. Первая регистрация
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "john_doe",
    "password": "password123"
  }'
```

### 4. Логин и получение JWT
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Полные примеры смотри в [EXAMPLES.md](./EXAMPLES.md) 👈

---

## 📊 Архитектура

```
┌─────────────────────────────────────────────┐
│         Frontend (React/Vue/Angular)        │
└──────────────┬──────────────────────────────┘
               │ HTTP Requests
               ▼
┌──────────────────────────────────────────────┐
│            Gin HTTP Server                    │
│ ┌────────────────────────────────────────┐  │
│ │  Handlers (auth, token management)     │  │
│ └────────────┬─────────────────────────┬─┘  │
│              │                         │     │
│  ┌───────────▼────────┐   ┌────────────▼──┐ │
│  │ AuthService        │   │ TokenService   │ │
│  │ (Register, Login)  │   │ (Generate API) │ │
│  └───────────┬────────┘   └────────────┬──┘ │
│              │                         │     │
│  ┌───────────▼─────────────────────────▼──┐ │
│  │   Repository Layer (Data Access)       │ │
│  └───────────┬──────────────────────────┬─┘ │
│              │                          │    │
└──────────────┼──────────────────────────┼────┘
               │                          │
        ┌──────▼──────┐          ┌────────▼────────┐
        │ PostgreSQL  │          │  SQLC Generated │
        │ Database    │◄─────────│  Code & Models  │
        └─────────────┘          └─────────────────┘
```

---

## 🔑 Аутентификация

### JWT Token (для основного API)
```
Authorization: Bearer <JWT_TOKEN>
```
- ⏰ Действует 24 часа
- 🔐 Выдается при успешном логине
- 📱 Используется для веб и мобильных приложений

### API Token (для консольных инструментов)
```
X-API-Token: <API_TOKEN>
```
- ♾️ Действует пока активна подписка
- 🔒 Хранится как SHA256 хеш
- 🛠️ Используется для автоматизации и скриптов

**Подробнее:** [API.md - Аутентификация](./API.md#аутентификация)

---

## 📦 Основные features

✅ **Регистрация и авторизация**
- Email/пароль аутентификация
- Bcrypt хеширование паролей
- JWT токены для сессий
- Аудит логирование

✅ **API Token Management**
- Генерация API токенов
- SHA256 хеширование токенов
- Максимум 10 активных токенов
- История использования токенов

✅ **Подписки**
- Trial, Premium, Commercial планы
- Автоматический trial при регистрации
- Лимиты на использование API
- Отслеживание использования

✅ **Безопасность**
- CORS поддержка (планируется)
- Rate limiting (планируется)
- Encryption for sensitive data
- Audit logs для всех действий

---

## 🌐 Endpoints (Quick Reference)

| Метод | Endpoint | Описание | Auth |
|-------|----------|---------|------|
| GET | `/health` | Проверка доступности | ❌ |
| POST | `/auth/register` | Регистрация | ❌ |
| POST | `/auth/login` | Авторизация | ❌ |
| POST | `/api/tokens` | Генерировать API токен | ✅ |
| GET | `/api/tokens` | Получить список токенов | ✅ |
| DELETE | `/api/tokens/:id` | Отозвать токен | ✅ |
| POST | `/api/logout` | Выход | ✅ |

**Полные детали:** [API.md - Endpoints](./API.md#endpoints)

---

## 🛠️ Технологический стек

### Backend
- 🐹 **Go 1.25** - Язык программирования
- 🔧 **Gin v1.12** - Web framework
- 🗄️ **PostgreSQL 18** - База данных
- 🔍 **SQLC v1.31** - Генерация типобезопасного кода
- 🔑 **pgx v5.9** - PostgreSQL драйвер
- 🎫 **golang-jwt v5** - JWT токены
- 🔐 **bcrypt** - Хеширование паролей
- 📋 **Goose** - Database migrations

### Database
- **Host**: localhost:5432
- **Database**: gogen
- **Driver**: PostgreSQL pgx v5

---

## 📋 Требования

- Go 1.24+
- PostgreSQL 16+
- Docker & Docker Compose (опционально)

---

## 🚦 Статус API

| Компонент | Статус | Версия |
|-----------|--------|--------|
| Core Auth | ✅ Готов | 1.0 |
| API Tokens | ✅ Готов | 1.0 |
| User Management | ✅ Готов | 1.0 |
| Subscriptions | ✅ Готов | 1.0 |
| Audit Logs | ✅ Готов | 1.0 |
| Email Verification | 🚧 В разработке | 1.1 |
| Password Reset | 🚧 В разработке | 1.1 |
| Rate Limiting | 🚧 Планируется | 1.1 |
| 2FA | 🚧 Планируется | 1.2 |

---

## 📞 Контакты и поддержка

- 📧 Email: support@example.com
- 🐛 Bug Reports: GitHub Issues
- 💡 Feature Requests: GitHub Discussions
- 📚 Documentation: Эта папка

---

## 📄 Лицензия

MIT License

---

## 🙏 Благодарности

Спасибо за использование GoGen API!

---

## 📖 Дополнительные ресурсы

- [Gin Documentation](https://gin-gonic.com/)
- [SQLC Documentation](https://sqlc.dev/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [JWT.io](https://jwt.io/)

---

**Последнее обновление**: 2026-05-18  
**Версия API**: 1.0  
**Статус**: Стабильная ✅
