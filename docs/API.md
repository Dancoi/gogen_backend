# API Documentation

## Базовый URL
```
http://localhost:8080
```

## Общая информация

### Статус коды ответов
- `200 OK` - Успешный запрос
- `201 Created` - Ресурс успешно создан
- `400 Bad Request` - Неверный формат запроса
- `401 Unauthorized` - Не авторизован
- `403 Forbidden` - Доступ запрещен
- `404 Not Found` - Ресурс не найден
- `409 Conflict` - Конфликт (например, email уже существует)
- `500 Internal Server Error` - Ошибка сервера

### Аутентификация

#### JWT Token (для основных API)
Все защищенные routes требуют заголовка:
```
Authorization: Bearer <JWT_TOKEN>
```

#### API Token (альтернативный способ)
```
X-API-Token: <API_TOKEN>
```

---

## Endpoints

### 1. Health Check

#### GET `/health`

Проверка доступности сервера.

**Параметры:** Нет

**Ответ:**
```json
{
  "status": "ok"
}
```

**Пример cURL:**
```bash
curl -X GET http://localhost:8080/health
```

---

## Auth Routes (Без авторизации)

### 2. Регистрация

#### POST `/auth/register`

Создает новый аккаунт пользователя.

**Body:**
```json
{
  "email": "user@example.com",
  "username": "john_doe",
  "password": "SecurePassword123"
}
```

**Валидация:**
- `email` - корректный email (обязательно)
- `username` - минимум 3, максимум 50 символов (обязательно)
- `password` - минимум 6 символов (обязательно)

**Успешный ответ (201 Created):**
```json
{
  "message": "user registered successfully",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "john_doe",
    "is_active": true,
    "created_at": "2026-05-18T10:30:00Z"
  }
}
```

**Ошибки:**
- `400 Bad Request` - Ошибка валидации
- `409 Conflict` - Email или username уже используется

**Пример cURL:**
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "john_doe",
    "password": "SecurePassword123"
  }'
```

---

### 3. Авторизация

#### POST `/auth/login`

Авторизует пользователя и выдает JWT токен.

**Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123"
}
```

**Валидация:**
- `email` - корректный email (обязательно)
- `password` - минимум 6 символов (обязательно)

**Успешный ответ (200 OK):**
```json
{
  "message": "login successful",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "john_doe",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Ошибки:**
- `400 Bad Request` - Ошибка валидации
- `401 Unauthorized` - Неверный email или пароль

**Пример cURL:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123"
  }'
```

---

## Protected Routes (Требует авторизации)

### 4. Генерировать API Токен

#### POST `/api/tokens`

Создает новый API токен для доступа через X-API-Token. **Максимум 10 активных токенов на пользователя.**

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

**Body:**
```json
{
  "name": "My Console App"
}
```

**Валидация:**
- `name` - минимум 1, максимум 255 символов (обязательно)

**Успешный ответ (201 Created):**
```json
{
  "message": "token generated successfully",
  "data": {
    "token": "gog_1234567890abcdef1234567890abcdef",
    "name": "My Console App",
    "note": "Save this token in a secure place. You won't be able to see it again!"
  }
}
```

**⚠️ Важно:** Токен выдается только один раз! Сохраните его в безопасном месте.

**Ошибки:**
- `400 Bad Request` - Ошибка валидации
- `401 Unauthorized` - Требуется авторизация
- `403 Forbidden` - Подписка истекла или достигнут лимит токенов

**Пример cURL:**
```bash
curl -X POST http://localhost:8080/api/tokens \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Console App"
  }'
```

---

### 5. Получить список API токенов

#### GET `/api/tokens`

Получает все активные API токены пользователя (отозванные токены не отображаются).

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
```

**Параметры:** Нет

**Успешный ответ (200 OK):**
```json
{
  "tokens": [
    {
      "id": 1,
      "name": "My Console App",
      "created_at": "2026-05-18T10:35:00Z",
      "last_used": "2026-05-18T10:36:15Z"
    },
    {
      "id": 2,
      "name": "Mobile App",
      "created_at": "2026-05-18T10:40:00Z",
      "last_used": ""
    }
  ]
}
```

**Ошибки:**
- `401 Unauthorized` - Требуется авторизация
- `500 Internal Server Error` - Ошибка сервера

**Пример cURL:**
```bash
curl -X GET http://localhost:8080/api/tokens \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

---

### 6. Отозвать API Токен

#### DELETE `/api/tokens/:id`

Отзывает (деактивирует) API токен. Отозванный токен больше не может быть использован.

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
```

**Параметры URL:**
- `id` - ID токена (число)

**Успешный ответ (200 OK):**
```json
{
  "message": "token revoked successfully"
}
```

**Ошибки:**
- `400 Bad Request` - Неверный формат ID
- `401 Unauthorized` - Требуется авторизация
- `404 Not Found` - Токен не найден
- `500 Internal Server Error` - Ошибка сервера

**Пример cURL:**
```bash
curl -X DELETE http://localhost:8080/api/tokens/1 \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

---

### 7. Выход (Logout)

#### POST `/api/logout`

Выходит из аккаунта, отзывая текущую сессию.

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
```

**Body:** Пусто

**Успешный ответ (200 OK):**
```json
{
  "message": "logout successful"
}
```

**Ошибки:**
- `401 Unauthorized` - Требуется авторизация
- `500 Internal Server Error` - Ошибка сервера

**Пример cURL:**
```bash
curl -X POST http://localhost:8080/api/logout \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

---

## Примеры использования

### Полный флоу регистрации → логина → генерация токена

#### 1. Регистрация
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "john_doe",
    "password": "SecurePassword123"
  }'
```

#### 2. Логин (получить JWT)
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123"
  }'

# Сохраняем полученный token в переменную
export JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

#### 3. Генерация API токена
```bash
curl -X POST http://localhost:8080/api/tokens \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "CLI Tool"
  }'

# Получаем ответ с токеном
# export API_TOKEN="gog_1234567890abcdef1234567890abcdef"
```

#### 4. Получить список токенов
```bash
curl -X GET http://localhost:8080/api/tokens \
  -H "Authorization: Bearer $JWT_TOKEN"
```

#### 5. Отозвать токен
```bash
curl -X DELETE http://localhost:8080/api/tokens/1 \
  -H "Authorization: Bearer $JWT_TOKEN"
```

#### 6. Выход
```bash
curl -X POST http://localhost:8080/api/logout \
  -H "Authorization: Bearer $JWT_TOKEN"
```

---

## Ограничения и Лимиты

### Лимиты API токенов
- **Максимум активных токенов:** 10 на пользователя
- **Время жизни токена:** 30 дней (для trial подписки)
- **Максимальное количество вызовов API:** 1000 в сутки (trial) / 10000 (premium) / 100000 (commercial)

### JWT Token
- **Время жизни:** 24 часа
- **Алгоритм:** HS256

### Валидация
- **Email:** Должен быть корректным email адресом
- **Username:** 3-50 символов, только буквы, цифры и подчеркивания
- **Password:** Минимум 6 символов
- **Token name:** 1-255 символов

---

## Обработка ошибок

Все ошибки возвращаются в формате:

```json
{
  "error": "описание ошибки"
}
```

### Типичные ошибки:

**Ошибка валидации:**
```json
{
  "error": "Key: 'RegisterRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```

**Email уже используется:**
```json
{
  "error": "email already exists"
}
```

**Неверные учетные данные:**
```json
{
  "error": "invalid email or password"
}
```

**Лимит токенов превышен:**
```json
{
  "error": "token generation limit exceeded"
}
```

**Подписка истекла:**
```json
{
  "error": "subscription has expired"
}
```

**Требуется авторизация:**
```json
{
  "error": "missing authorization header"
}
```

---

## Headers

### Запросы
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>  (для защищенных routes)
```

### Ответы
```
Content-Type: application/json
```

---

## Версия API
- **Версия:** 1.0
- **Последнее обновление:** 2026-05-18
- **Статус:** Стабильная
