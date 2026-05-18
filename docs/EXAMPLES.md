# Примеры использования API

## TypeScript / React

### 1. Регистрация

```typescript
const register = async (email: string, username: string, password: string) => {
  const response = await fetch('http://localhost:8080/auth/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, username, password }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error);
  }

  const data = await response.json();
  return data.user;
};
```

### 2. Логин и сохранение JWT

```typescript
const login = async (email: string, password: string) => {
  const response = await fetch('http://localhost:8080/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error);
  }

  const data = await response.json();
  
  // Сохраняем JWT токен в localStorage
  localStorage.setItem('jwtToken', data.user.token);
  
  return data.user;
};
```

### 3. Генерирование API токена

```typescript
const generateApiToken = async (tokenName: string) => {
  const jwtToken = localStorage.getItem('jwtToken');
  
  const response = await fetch('http://localhost:8080/api/tokens', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${jwtToken}`,
    },
    body: JSON.stringify({ name: tokenName }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error);
  }

  const data = await response.json();
  
  // ВАЖНО: Показать пользователю токен (выводится только один раз!)
  alert(`Ваш API токен: ${data.data.token}\n\n${data.data.note}`);
  
  return data.data;
};
```

### 4. Получить список токенов

```typescript
const getApiTokens = async () => {
  const jwtToken = localStorage.getItem('jwtToken');
  
  const response = await fetch('http://localhost:8080/api/tokens', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${jwtToken}`,
    },
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error);
  }

  const data = await response.json();
  return data.tokens;
};
```

### 5. Отозвать API токен

```typescript
const revokeApiToken = async (tokenId: number) => {
  const jwtToken = localStorage.getItem('jwtToken');
  
  const response = await fetch(`http://localhost:8080/api/tokens/${tokenId}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${jwtToken}`,
    },
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error);
  }

  const data = await response.json();
  return data;
};
```

### 6. Выход из аккаунта

```typescript
const logout = async () => {
  const jwtToken = localStorage.getItem('jwtToken');
  
  const response = await fetch('http://localhost:8080/api/logout', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${jwtToken}`,
    },
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error);
  }

  // Удаляем токен из localStorage
  localStorage.removeItem('jwtToken');
};
```

### 7. Axios Helper

```typescript
import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080';

// Создаем axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
});

// Добавляем JWT токен в каждый запрос
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('jwtToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Экспортируем функции
export const authAPI = {
  register: (email: string, username: string, password: string) =>
    api.post('/auth/register', { email, username, password }),
  
  login: (email: string, password: string) =>
    api.post('/auth/login', { email, password }),
  
  logout: () =>
    api.post('/api/logout'),
};

export const tokensAPI = {
  generate: (name: string) =>
    api.post('/api/tokens', { name }),
  
  list: () =>
    api.get('/api/tokens'),
  
  revoke: (id: number) =>
    api.delete(`/api/tokens/${id}`),
};
```

---

## JavaScript (Vanilla)

### Простой HTTP Client

```javascript
class ApiClient {
  constructor(baseURL = 'http://localhost:8080') {
    this.baseURL = baseURL;
  }

  async request(endpoint, options = {}) {
    const token = localStorage.getItem('jwtToken');
    const headers = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${this.baseURL}${endpoint}`, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'API Error');
    }

    return await response.json();
  }

  // Auth endpoints
  register(email, username, password) {
    return this.request('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, username, password }),
    });
  }

  login(email, password) {
    return this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
  }

  logout() {
    return this.request('/api/logout', { method: 'POST' });
  }

  // Token endpoints
  generateToken(name) {
    return this.request('/api/tokens', {
      method: 'POST',
      body: JSON.stringify({ name }),
    });
  }

  listTokens() {
    return this.request('/api/tokens', { method: 'GET' });
  }

  revokeToken(id) {
    return this.request(`/api/tokens/${id}`, { method: 'DELETE' });
  }
}

// Использование
const api = new ApiClient();

// Логин
api.login('user@example.com', 'password123')
  .then(data => {
    localStorage.setItem('jwtToken', data.user.token);
    console.log('Успешный вход:', data.user);
  })
  .catch(error => console.error('Ошибка:', error.message));
```

---

## Python

### Requests Library

```python
import requests
import json

BASE_URL = 'http://localhost:8080'

class GoGenAPI:
    def __init__(self, base_url=BASE_URL):
        self.base_url = base_url
        self.session = requests.Session()
        self.token = None
    
    def set_token(self, token):
        self.token = token
        self.session.headers.update({'Authorization': f'Bearer {token}'})
    
    def register(self, email, username, password):
        response = self.session.post(
            f'{self.base_url}/auth/register',
            json={'email': email, 'username': username, 'password': password}
        )
        return response.json()
    
    def login(self, email, password):
        response = self.session.post(
            f'{self.base_url}/auth/login',
            json={'email': email, 'password': password}
        )
        data = response.json()
        self.set_token(data['user']['token'])
        return data
    
    def logout(self):
        response = self.session.post(f'{self.base_url}/api/logout')
        self.session.headers.pop('Authorization', None)
        return response.json()
    
    def generate_token(self, name):
        response = self.session.post(
            f'{self.base_url}/api/tokens',
            json={'name': name}
        )
        return response.json()
    
    def list_tokens(self):
        response = self.session.get(f'{self.base_url}/api/tokens')
        return response.json()
    
    def revoke_token(self, token_id):
        response = self.session.delete(f'{self.base_url}/api/tokens/{token_id}')
        return response.json()

# Использование
api = GoGenAPI()

# Логин
login_data = api.login('user@example.com', 'password123')
print('User:', login_data['user'])

# Генерирование токена
token_data = api.generate_token('CLI Tool')
print('New API Token:', token_data['data']['token'])

# Список токенов
tokens = api.list_tokens()
print('Tokens:', tokens['tokens'])

# Отозвать токен
api.revoke_token(1)
print('Token revoked')

# Выход
api.logout()
```

---

## cURL (Bash)

### Полный пример

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

# 1. Регистрация
echo "=== Регистрация ==="
REGISTER=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "john_doe",
    "password": "SecurePass123"
  }')
echo $REGISTER

# 2. Логин
echo -e "\n=== Логин ==="
LOGIN=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123"
  }')
echo $LOGIN

# Извлекаем JWT токен
JWT_TOKEN=$(echo $LOGIN | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "JWT Token: $JWT_TOKEN"

# 3. Генерирование API токена
echo -e "\n=== Генерирование API токена ==="
TOKEN=$(curl -s -X POST "$BASE_URL/api/tokens" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"CLI Tool"}')
echo $TOKEN

# 4. Получить список токенов
echo -e "\n=== Список токенов ==="
curl -s -X GET "$BASE_URL/api/tokens" \
  -H "Authorization: Bearer $JWT_TOKEN"

# 5. Выход
echo -e "\n=== Выход ==="
curl -s -X POST "$BASE_URL/api/logout" \
  -H "Authorization: Bearer $JWT_TOKEN"
```

---

## Go

### Using net/http

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const BaseURL = "http://localhost:8080"

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func NewClient() *Client {
	return &Client{
		baseURL: BaseURL,
		http:    &http.Client{},
	}
}

func (c *Client) do(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, _ := http.NewRequest(method, c.baseURL+endpoint, reqBody)
	req.Header.Set("Content-Type", "application/json")

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *Client) Login(email, password string) error {
	resp, _ := c.do("POST", "/auth/login", map[string]string{
		"email":    email,
		"password": password,
	})

	var data map[string]interface{}
	json.Unmarshal(resp, &data)
	user := data["user"].(map[string]interface{})
	c.token = user["token"].(string)
	return nil
}

func (c *Client) GenerateToken(name string) (string, error) {
	resp, _ := c.do("POST", "/api/tokens", map[string]string{"name": name})

	var data map[string]interface{}
	json.Unmarshal(resp, &data)
	tokenData := data["data"].(map[string]interface{})
	return tokenData["token"].(string), nil
}

// Использование
func main() {
	client := NewClient()
	client.Login("user@example.com", "password123")
	token, _ := client.GenerateToken("My CLI Tool")
	fmt.Println("API Token:", token)
}
```

---

## Postman Collection

Импортируйте этот JSON в Postman для быстрого тестирования:

```json
{
  "info": {
    "name": "GoGen API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Auth",
      "item": [
        {
          "name": "Register",
          "request": {
            "method": "POST",
            "url": "http://localhost:8080/auth/register",
            "header": [
              {"key": "Content-Type", "value": "application/json"}
            ],
            "body": {
              "mode": "raw",
              "raw": "{\"email\":\"user@example.com\",\"username\":\"john_doe\",\"password\":\"SecurePass123\"}"
            }
          }
        },
        {
          "name": "Login",
          "request": {
            "method": "POST",
            "url": "http://localhost:8080/auth/login",
            "header": [
              {"key": "Content-Type", "value": "application/json"}
            ],
            "body": {
              "mode": "raw",
              "raw": "{\"email\":\"user@example.com\",\"password\":\"SecurePass123\"}"
            }
          }
        }
      ]
    },
    {
      "name": "Tokens",
      "item": [
        {
          "name": "Generate Token",
          "request": {
            "method": "POST",
            "url": "http://localhost:8080/api/tokens",
            "header": [
              {"key": "Authorization", "value": "Bearer {{jwt_token}}"},
              {"key": "Content-Type", "value": "application/json"}
            ],
            "body": {
              "mode": "raw",
              "raw": "{\"name\":\"My CLI Tool\"}"
            }
          }
        },
        {
          "name": "List Tokens",
          "request": {
            "method": "GET",
            "url": "http://localhost:8080/api/tokens",
            "header": [
              {"key": "Authorization", "value": "Bearer {{jwt_token}}"}
            ]
          }
        }
      ]
    }
  ]
}
```
