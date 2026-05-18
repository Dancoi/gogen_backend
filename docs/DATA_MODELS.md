# Структуры данных

## User (Пользователь)

```json
{
  "id": 1,
  "email": "user@example.com",
  "username": "john_doe",
  "password_hash": "$2a$10$...",
  "is_active": true,
  "email_verified": false,
  "created_at": "2026-05-18T10:30:00Z",
  "updated_at": "2026-05-18T10:30:00Z"
}
```

### Поля:
- **id** (int32) - Уникальный идентификатор
- **email** (string) - Email адрес (уникальный)
- **username** (string) - Имя пользователя (уникальное)
- **password_hash** (string) - Хеш пароля (bcrypt)
- **is_active** (bool) - Активен ли пользователь
- **email_verified** (bool) - Подтвержден ли email
- **created_at** (timestamp) - Время создания
- **updated_at** (timestamp) - Время последнего обновления

---

## ApiToken (API Токен)

```json
{
  "id": 1,
  "user_id": 1,
  "subscription_id": 5,
  "name": "My CLI Tool",
  "token_hash": "sha256_hash_of_token",
  "is_active": true,
  "created_at": "2026-05-18T10:35:00Z",
  "last_used_at": "2026-05-18T10:36:15Z",
  "expires_at": "2026-06-18T10:35:00Z",
  "revoked_at": null
}
```

### Поля:
- **id** (int32) - Уникальный идентификатор
- **user_id** (int32) - ID пользователя (внешний ключ)
- **subscription_id** (int32) - ID подписки (внешний ключ)
- **name** (string) - Имя/описание токена
- **token_hash** (string) - SHA256 хеш токена (исходный токен не сохраняется!)
- **is_active** (bool) - Активен ли токен
- **created_at** (timestamp) - Время создания
- **last_used_at** (timestamp) - Время последнего использования
- **expires_at** (timestamp) - Время истечения
- **revoked_at** (timestamp, nullable) - Время отзыва (если отозван)

### Примечания:
- ⚠️ **Исходный токен выдается только один раз!** Он не сохраняется в БД
- В БД хранится только SHA256 хеш токена
- Максимум 10 активных токенов на пользователя
- Токен может использоваться через заголовок `X-API-Token`

---

## UserSession (Сессия пользователя)

```json
{
  "id": 1,
  "user_id": 1,
  "session_token": "random_token_string",
  "ip_address": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "expires_at": "2026-05-19T10:30:00Z",
  "last_activity_at": "2026-05-18T10:35:00Z",
  "created_at": "2026-05-18T10:30:00Z"
}
```

### Поля:
- **id** (int32) - Уникальный идентификатор
- **user_id** (int32) - ID пользователя (внешний ключ)
- **session_token** (string) - Токен сессии (уникальный)
- **ip_address** (string, nullable) - IP адрес подключения
- **user_agent** (string, nullable) - User-Agent браузера
- **expires_at** (timestamp) - Время истечения сессии
- **last_activity_at** (timestamp) - Время последней активности
- **created_at** (timestamp) - Время создания

---

## Subscription (Подписка)

```json
{
  "id": 1,
  "user_id": 1,
  "subscription_plan_id": 1,
  "plan_type": "trial",
  "is_active": true,
  "started_at": "2026-05-18T10:30:00Z",
  "expires_at": "2026-06-18T10:30:00Z",
  "renewal_date": null,
  "max_tokens_per_month": 100,
  "current_usage": 25,
  "usage_reset_date": "2026-06-01T00:00:00Z",
  "created_at": "2026-05-18T10:30:00Z",
  "updated_at": "2026-05-18T10:30:00Z"
}
```

### Поля:
- **id** (int32) - Уникальный идентификатор
- **user_id** (int32) - ID пользователя (внешний ключ)
- **subscription_plan_id** (int32) - ID плана подписки (внешний ключ)
- **plan_type** (string) - Тип плана: "trial", "premium", "commercial"
- **is_active** (bool) - Активна ли подписка
- **started_at** (timestamp) - Дата начала
- **expires_at** (timestamp) - Дата истечения
- **renewal_date** (timestamp, nullable) - Дата автоматического продления
- **max_tokens_per_month** (int32) - Максимум токенов в месяц
- **current_usage** (int32) - Текущее использование
- **usage_reset_date** (timestamp) - Дата сброса счетчика
- **created_at** (timestamp) - Время создания
- **updated_at** (timestamp) - Время последнего обновления

---

## SubscriptionPlan (План подписки)

```json
{
  "id": 1,
  "name": "Trial Plan",
  "plan_type": "trial",
  "max_tokens_per_month": 100,
  "max_api_calls_per_day": 1000,
  "price": 0.00,
  "features": {
    "console_tool": true,
    "api_access": true,
    "support": "email"
  },
  "trial_duration_days": 30,
  "is_active": true,
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-01T00:00:00Z"
}
```

### Поля:
- **id** (int32) - Уникальный идентификатор
- **name** (string) - Название плана
- **plan_type** (string) - Тип плана (уникальный): "trial", "premium", "commercial"
- **max_tokens_per_month** (int32) - Максимум API токенов в месяц
- **max_api_calls_per_day** (int32) - Максимум API вызовов в сутки
- **price** (decimal) - Цена (в USD)
- **features** (jsonb) - Функции плана
- **trial_duration_days** (int32) - Длительность trial периода в днях
- **is_active** (bool) - Активен ли план
- **created_at** (timestamp) - Время создания
- **updated_at** (timestamp) - Время последнего обновления

### Предустановленные планы:
```json
[
  {
    "id": 1,
    "name": "Trial",
    "plan_type": "trial",
    "max_tokens_per_month": 100,
    "max_api_calls_per_day": 1000,
    "price": 0.00,
    "trial_duration_days": 30
  },
  {
    "id": 2,
    "name": "Premium",
    "plan_type": "premium",
    "max_tokens_per_month": 1000,
    "max_api_calls_per_day": 10000,
    "price": 29.99
  },
  {
    "id": 3,
    "name": "Commercial",
    "plan_type": "commercial",
    "max_tokens_per_month": 10000,
    "max_api_calls_per_day": 100000,
    "price": 299.99
  }
]
```

---

## AuditLog (Логи аудита)

```json
{
  "id": 1,
  "user_id": 1,
  "action": "user_created",
  "resource_type": "user",
  "resource_id": 1,
  "changes": {
    "email": "user@example.com",
    "username": "john_doe"
  },
  "status_code": 201,
  "error_message": null,
  "ip_address": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "created_at": "2026-05-18T10:30:00Z"
}
```

### Поля:
- **id** (int32) - Уникальный идентификатор
- **user_id** (int32, nullable) - ID пользователя
- **action** (string) - Действие: "user_created", "login", "token_created", "token_revoked" и т.д.
- **resource_type** (string, nullable) - Тип ресурса: "user", "token", "subscription"
- **resource_id** (int32, nullable) - ID ресурса
- **changes** (jsonb) - Данные об изменениях
- **status_code** (int32, nullable) - HTTP статус код
- **error_message** (string, nullable) - Сообщение об ошибке (если была)
- **ip_address** (string, nullable) - IP адрес
- **user_agent** (string, nullable) - User-Agent
- **created_at** (timestamp) - Время создания

---

## Role (Роль пользователя)

```json
{
  "id": 1,
  "name": "admin",
  "description": "Administrator with full access",
  "permissions": [
    "users:read",
    "users:write",
    "users:delete",
    "subscriptions:read",
    "subscriptions:write"
  ],
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-01T00:00:00Z"
}
```

### Поля:
- **id** (int32) - Уникальный идентификатор
- **name** (string) - Название роли (уникальное)
- **description** (string) - Описание
- **permissions** (jsonb) - Список разрешений
- **created_at** (timestamp) - Время создания
- **updated_at** (timestamp) - Время последнего обновления

---

## TokenUsageLog (Логи использования токенов)

```json
{
  "id": 1,
  "api_token_id": 1,
  "user_id": 1,
  "endpoint": "/api/generate",
  "method": "POST",
  "status_code": 200,
  "response_time": 125,
  "data_size": 2048,
  "error_message": null,
  "ip_address": "192.168.1.100",
  "user_agent": "CLI Tool v1.0",
  "created_at": "2026-05-18T10:35:00Z"
}
```

### Поля:
- **id** (int64) - Уникальный идентификатор
- **api_token_id** (int32) - ID API токена (внешний ключ)
- **user_id** (int32) - ID пользователя (внешний ключ)
- **endpoint** (string) - Endpoint который был вызван
- **method** (string) - HTTP метод: GET, POST, PUT, DELETE
- **status_code** (int32) - HTTP статус код ответа
- **response_time** (int64) - Время ответа в миллисекундах
- **data_size** (int64) - Размер данных в байтах
- **error_message** (string, nullable) - Сообщение об ошибке (если была)
- **ip_address** (string, nullable) - IP адрес клиента
- **user_agent** (string, nullable) - User-Agent
- **created_at** (timestamp) - Время логирования

---

## Ответ при ошибке

```json
{
  "error": "описание ошибки"
}
```

### Примеры ошибок:

**Email уже используется:**
```json
{"error": "email already exists"}
```

**Username уже используется:**
```json
{"error": "username already exists"}
```

**Неверные учетные данные:**
```json
{"error": "invalid email or password"}
```

**Ошибка валидации:**
```json
{"error": "Key: 'RegisterRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"}
```

**Лимит токенов:**
```json
{"error": "token generation limit exceeded"}
```

**Подписка истекла:**
```json
{"error": "subscription has expired"}
```

**Требуется авторизация:**
```json
{"error": "missing authorization header"}
```

**Неверный формат авторизации:**
```json
{"error": "invalid authorization header format"}
```

**Токен не найден:**
```json
{"error": "token not found"}
```

---

## Типы данных

### Временные метки (Timestamp)
Все временные метки в формате ISO 8601:
```
2026-05-18T10:30:00Z
```

### Булевы значения (Bool)
```
true / false
```

### Числовые значения
- **int32** - Целое число (−2,147,483,648 до 2,147,483,647)
- **int64** - Большое целое число (−9,223,372,036,854,775,808 до 9,223,372,036,854,775,807)
- **decimal** - Число с плавающей точкой (для цен)

### Строки (String)
UTF-8 кодировка

### JSON объекты
Вложенные JSON структуры для гибких данных (features, permissions, changes)

---

## Связи между таблицами

```
User (1) ──── (N) Subscription
User (1) ──── (N) ApiToken
User (1) ──── (N) UserSession
User (1) ──── (N) AuditLog
User (N) ──── (N) Role (через user_roles junction table)

SubscriptionPlan (1) ──── (N) Subscription

Subscription (1) ──── (N) ApiToken

ApiToken (1) ──── (N) TokenUsageLog
```
