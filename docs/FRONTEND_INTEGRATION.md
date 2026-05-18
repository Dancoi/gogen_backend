# Интеграция API с фронтенд приложением 🚀

Пошаговая инструкция по интеграции GoGen API с вашим фронтенд приложением.

## Оглавление
1. [Подготовка](#подготовка)
2. [Настройка окружения](#настройка-окружения)
3. [API Client](#api-client)
4. [Реализация функций](#реализация-функций)
5. [Обработка ошибок](#обработка-ошибок)
6. [Тестирование](#тестирование)

---

## Подготовка

### 1. Проверить что backend запущен
```bash
# В отдельном терминале
cd gogen_backend
go run ./cmd/main/main.go

# Проверить доступность
curl http://localhost:8080/health
# Должен вернуть: {"status":"ok"}
```

### 2. Установить необходимые зависимости

#### Для React
```bash
npm install axios
# или fetch встроенный (не требует установки)
```

#### Для Vue
```bash
npm install axios
```

#### Для Angular
```bash
# HttpClient уже встроен
```

---

## Настройка окружения

### 1. Создать `.env` файл в корне фронтенд проекта

#### React (.env)
```env
REACT_APP_API_URL=http://localhost:8080
REACT_APP_API_TIMEOUT=5000
```

#### Vue (.env.local)
```env
VUE_APP_API_URL=http://localhost:8080
VUE_APP_API_TIMEOUT=5000
```

#### Angular (environment.ts)
```typescript
export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080',
  apiTimeout: 5000,
};
```

### 2. Использовать переменные окружения

#### React
```typescript
const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
```

#### Vue
```typescript
const API_URL = process.env.VUE_APP_API_URL || 'http://localhost:8080';
```

#### Angular
```typescript
import { environment } from '../environments/environment';
const API_URL = environment.apiUrl;
```

---

## API Client

### React / Vanilla JS - Axios

#### services/api.ts
```typescript
import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_URL,
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Добавляем JWT токен в каждый запрос
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('jwtToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Обработка ошибок
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Токен истек или невалиден
      localStorage.removeItem('jwtToken');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api;
```

#### services/auth.ts
```typescript
import api from './api';

export const authService = {
  register: (email: string, username: string, password: string) =>
    api.post('/auth/register', { email, username, password }),

  login: (email: string, password: string) =>
    api.post('/auth/login', { email, password }),

  logout: () => {
    api.post('/api/logout');
    localStorage.removeItem('jwtToken');
  },
};
```

#### services/tokens.ts
```typescript
import api from './api';

export const tokensService = {
  generate: (name: string) =>
    api.post('/api/tokens', { name }),

  list: () =>
    api.get('/api/tokens'),

  revoke: (id: number) =>
    api.delete(`/api/tokens/${id}`),
};
```

---

## Реализация функций

### React - Complete Example

#### hooks/useAuth.ts
```typescript
import { useState } from 'react';
import { authService } from '../services/auth';

export const useAuth = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [user, setUser] = useState(null);

  const login = async (email: string, password: string) => {
    setLoading(true);
    setError(null);
    try {
      const response = await authService.login(email, password);
      const { token, ...userData } = response.data.user;
      
      // Сохраняем JWT
      localStorage.setItem('jwtToken', token);
      setUser(userData);
      
      return userData;
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Login failed';
      setError(errorMsg);
      throw new Error(errorMsg);
    } finally {
      setLoading(false);
    }
  };

  const register = async (email: string, username: string, password: string) => {
    setLoading(true);
    setError(null);
    try {
      const response = await authService.register(email, username, password);
      setUser(response.data.user);
      return response.data.user;
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Registration failed';
      setError(errorMsg);
      throw new Error(errorMsg);
    } finally {
      setLoading(false);
    }
  };

  const logout = () => {
    authService.logout();
    setUser(null);
  };

  return { user, loading, error, login, register, logout };
};
```

#### components/LoginForm.tsx
```typescript
import React, { useState } from 'react';
import { useAuth } from '../hooks/useAuth';

export const LoginForm = () => {
  const { login, loading, error } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await login(email, password);
      // Перенаправить на dashboard
      window.location.href = '/dashboard';
    } catch (err) {
      // Ошибка уже установлена в hook
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      {error && <div className="error">{error}</div>}
      
      <input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder="Email"
        required
      />
      
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Password"
        required
      />
      
      <button type="submit" disabled={loading}>
        {loading ? 'Logging in...' : 'Login'}
      </button>
    </form>
  );
};
```

#### components/TokenGenerator.tsx
```typescript
import React, { useState, useEffect } from 'react';
import { tokensService } from '../services/tokens';

export const TokenGenerator = () => {
  const [tokens, setTokens] = useState([]);
  const [tokenName, setTokenName] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadTokens();
  }, []);

  const loadTokens = async () => {
    try {
      const response = await tokensService.list();
      setTokens(response.data.tokens);
    } catch (err) {
      setError('Failed to load tokens');
    }
  };

  const handleGenerateToken = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const response = await tokensService.generate(tokenName);
      const newToken = response.data.data.token;
      
      // Показать токен пользователю
      alert(`Your API Token:\n${newToken}\n\n${response.data.data.note}`);
      
      // Перезагрузить список
      setTokenName('');
      await loadTokens();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to generate token');
    } finally {
      setLoading(false);
    }
  };

  const handleRevokeToken = async (tokenId: number) => {
    if (window.confirm('Are you sure you want to revoke this token?')) {
      try {
        await tokensService.revoke(tokenId);
        await loadTokens();
      } catch (err) {
        setError('Failed to revoke token');
      }
    }
  };

  return (
    <div>
      <h2>API Tokens</h2>

      {error && <div className="error">{error}</div>}

      <form onSubmit={handleGenerateToken}>
        <input
          type="text"
          value={tokenName}
          onChange={(e) => setTokenName(e.target.value)}
          placeholder="Token name (e.g., 'CLI Tool')"
          maxLength={255}
          required
        />
        <button type="submit" disabled={loading}>
          {loading ? 'Generating...' : 'Generate Token'}
        </button>
      </form>

      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Created</th>
            <th>Last Used</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          {tokens.map((token) => (
            <tr key={token.id}>
              <td>{token.name}</td>
              <td>{new Date(token.created_at).toLocaleDateString()}</td>
              <td>{token.last_used || 'Never'}</td>
              <td>
                <button onClick={() => handleRevokeToken(token.id)}>
                  Revoke
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
```

---

### Vue 3 - Complete Example

#### services/api.ts (для Vue)
```typescript
import axios from 'axios';

const API_URL = process.env.VUE_APP_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_URL,
  timeout: 5000,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('jwtToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('jwtToken');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api;
```

#### composables/useAuth.ts
```typescript
import { ref } from 'vue';
import api from '../services/api';

export const useAuth = () => {
  const user = ref(null);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const login = async (email: string, password: string) => {
    loading.value = true;
    error.value = null;
    try {
      const response = await api.post('/auth/login', { email, password });
      const { token, ...userData } = response.data.user;
      localStorage.setItem('jwtToken', token);
      user.value = userData;
      return userData;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Login failed';
      throw error.value;
    } finally {
      loading.value = false;
    }
  };

  const register = async (email: string, username: string, password: string) => {
    loading.value = true;
    error.value = null;
    try {
      const response = await api.post('/auth/register', { email, username, password });
      user.value = response.data.user;
      return response.data.user;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Registration failed';
      throw error.value;
    } finally {
      loading.value = false;
    }
  };

  const logout = () => {
    api.post('/api/logout');
    localStorage.removeItem('jwtToken');
    user.value = null;
  };

  return { user, loading, error, login, register, logout };
};
```

#### components/LoginForm.vue
```vue
<template>
  <form @submit.prevent="handleSubmit">
    <div v-if="error" class="error">{{ error }}</div>
    
    <input
      v-model="email"
      type="email"
      placeholder="Email"
      required
    />
    
    <input
      v-model="password"
      type="password"
      placeholder="Password"
      required
    />
    
    <button :disabled="loading">
      {{ loading ? 'Logging in...' : 'Login' }}
    </button>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useAuth } from '../composables/useAuth';
import { useRouter } from 'vue-router';

const { login, loading, error } = useAuth();
const router = useRouter();
const email = ref('');
const password = ref('');

const handleSubmit = async () => {
  try {
    await login(email.value, password.value);
    router.push('/dashboard');
  } catch (err) {
    // Ошибка уже установлена
  }
};
</script>
```

---

## Обработка ошибок

### Типичные ошибки API и как их обрабатывать

```typescript
const handleApiError = (error: any) => {
  if (!error.response) {
    // Сетевая ошибка
    return 'Network error. Please check your connection.';
  }

  const status = error.response.status;
  const data = error.response.data;

  switch (status) {
    case 400:
      return `Invalid request: ${data.error}`;
    case 401:
      return 'Unauthorized. Please login again.';
    case 403:
      return `Access forbidden: ${data.error}`;
    case 404:
      return 'Resource not found.';
    case 409:
      return `Conflict: ${data.error}`;
    case 500:
      return 'Server error. Please try again later.';
    default:
      return `Error: ${data.error || 'Something went wrong'}`;
  }
};
```

### Error Boundary для React

```typescript
import React from 'react';

class ErrorBoundary extends React.Component<any, any> {
  constructor(props: any) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: any) {
    return { hasError: true, error };
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="error">
          <h1>Something went wrong</h1>
          <p>{this.state.error?.message}</p>
        </div>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
```

---

## Тестирование

### Тест регистрации и логина

```typescript
describe('Authentication', () => {
  it('should register a new user', async () => {
    const response = await authService.register(
      'test@example.com',
      'testuser',
      'password123'
    );
    expect(response.data.user.email).toBe('test@example.com');
  });

  it('should login with valid credentials', async () => {
    const response = await authService.login(
      'test@example.com',
      'password123'
    );
    expect(response.data.user.token).toBeDefined();
    expect(localStorage.getItem('jwtToken')).toBe(response.data.user.token);
  });

  it('should fail login with invalid credentials', async () => {
    try {
      await authService.login('test@example.com', 'wrongpassword');
      fail('Should throw error');
    } catch (error: any) {
      expect(error.response.status).toBe(401);
    }
  });
});
```

---

## Production Deployment

### Environment Variables для Production

#### .env.production
```env
REACT_APP_API_URL=https://api.yourdomain.com
REACT_APP_API_TIMEOUT=5000
```

### CORS Configuration

На production backend должен разрешить CORS для вашего домена:

```go
// Будет добавлено в main.go
config := cors.DefaultConfig()
config.AllowOrigins = []string{"https://yourdomain.com"}
router.Use(cors.New(config))
```

### Build для Production

#### React
```bash
npm run build
# Результат в папке build/
```

#### Vue
```bash
npm run build
# Результат в папке dist/
```

---

## Checklist интеграции

- [ ] API Client создан и настроен
- [ ] JWT токен сохраняется в localStorage
- [ ] JWT токен отправляется в каждом запросе
- [ ] Ошибки API обрабатываются правильно
- [ ] Требуется авторизация для защищенных routes
- [ ] При истечении JWT пользователь перенаправляется на login
- [ ] API токены отображаются правильно
- [ ] Можно генерировать новые API токены
- [ ] Можно отозвать API токены
- [ ] UI обновляется после каждого действия
- [ ] Локально протестировано с backend
- [ ] Production переменные окружения настроены
- [ ] CORS настроены для production
- [ ] Документация обновлена для team

---

## Полезные ссылки

- [Полная документация API](./API.md)
- [Примеры кода](./EXAMPLES.md)
- [Структуры данных](./DATA_MODELS.md)
- [FAQ](./FAQ.md)

---

**Успешной интеграции! 🚀**
