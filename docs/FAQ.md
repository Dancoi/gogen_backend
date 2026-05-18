# FAQ - Часто задаваемые вопросы

## Аутентификация и авторизация

### Q: Как сохранить JWT токен?
**A:** JWT токен выдается при успешном логине. Сохраняйте его в localStorage (для веб) или в защищенном хранилище (для мобильных приложений):

```javascript
// localStorage
localStorage.setItem('jwtToken', token);

// Получение
const token = localStorage.getItem('jwtToken');

// Удаление
localStorage.removeItem('jwtToken');
```

### Q: Как долго действует JWT токен?
**A:** JWT токен действует 24 часа. После истечения необходимо авторизоваться заново.

### Q: Можно ли использовать JWT токен после логаута?
**A:** Технически токен будет оставаться валидным до истечения 24 часов, но рекомендуется удалить его с клиента при логауте.

### Q: Как использовать API токен вместо JWT?
**A:** Отправьте API токен в заголовке `X-API-Token`:

```bash
curl -X GET http://localhost:8080/api/tokens \
  -H "X-API-Token: gog_1234567890abcdef1234567890abcdef"
```

---

## API Токены

### Q: Почему мой API токен не сохраняется?
**A:** API токен выдается ТОЛЬКО один раз! Если вы потеряли токен:
1. Отзовите старый токен через `/api/tokens/:id`
2. Сгенерируйте новый токен через `POST /api/tokens`
3. Сразу же сохраните новый токен

### Q: Сколько API токенов я могу создать?
**A:** Максимум 10 активных токенов на пользователя. Если достигнут лимит:
1. Отзовите неиспользуемые токены
2. Попробуйте создать новый

### Q: Как использовать API токен?
**A:** Используйте заголовок `X-API-Token`:

```bash
curl -X GET http://localhost:8080/api/tokens \
  -H "X-API-Token: gog_1234567890abcdef"
```

Или в коде:

```javascript
const response = await fetch('http://localhost:8080/api/tokens', {
  headers: {
    'X-API-Token': 'gog_1234567890abcdef'
  }
});
```

### Q: Могу ли я восстановить отозванный токен?
**A:** Нет, отозванные токены невозможно восстановить. Создайте новый.

### Q: Когда нужно использовать API токены вместо JWT?
**A:** 
- **JWT**: Для веб-приложений, мобильных приложений с пользовательским интерфейсом
- **API токены**: Для консольных инструментов, скриптов, микросервисов, автоматизации

---

## Подписки и лимиты

### Q: Какие лимиты на trial подписке?
**A:** Trial подписка:
- **Максимум API токенов в месяц**: 100
- **Максимум API вызовов в сутки**: 1000
- **Длительность**: 30 дней
- **Стоимость**: Бесплатно

### Q: Что происходит когда истекает trial подписка?
**A:** Когда истекает trial подписка:
1. Пользователь не может генерировать новые API токены
2. Существующие токены прекращают работать (возвращают 403 Forbidden)
3. Пользователю предлагается обновить подписку

### Q: Как обновить подписку?
**A:** На данный момент обновление подписки требует связи с поддержкой. В будущих версиях будет доступно через API.

### Q: Как узнать когда истекает моя подписка?
**A:** Получите информацию о подписке (в разработке):
```bash
GET /api/subscription  # Планируется
```

---

## Ошибки

### Q: Получаю ошибку "email already exists"
**A:** Email уже зарегистрирован в системе. Используйте другой email или выполните вход если это ваш аккаунт.

### Q: Получаю ошибку "invalid email or password"
**A:** Проверьте что:
1. Email введен правильно (без пробелов)
2. Пароль введен правильно (учитываются прописные/строчные буквы)
3. Аккаунт существует в системе

### Q: Получаю "missing authorization header"
**A:** Убедитесь что вы отправляете заголовок Authorization для защищенных routes:

```
Authorization: Bearer <JWT_TOKEN>
```

### Q: Получаю "invalid authorization header format"
**A:** Формат заголовка должен быть: `Bearer <TOKEN>` (с пробелом между словом "Bearer" и самим токеном).

### Q: Получаю "token generation limit exceeded"
**A:** Вы достигли лимита активных токенов (10). Отзовите неиспользуемые токены:

```bash
curl -X DELETE http://localhost:8080/api/tokens/1 \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### Q: Получаю "subscription has expired"
**A:** Ваша подписка истекла. Обновите подписку чтобы продолжить использовать API.

---

## Безопасность

### Q: Безопасно ли сохранять JWT в localStorage?
**A:** localStorage уязвим для XSS атак. Лучше использовать:
- HttpOnly cookies (автоматически отправляются в каждом запросе)
- Защищенное хранилище в мобильных приложениях

Однако для простых приложений localStorage приемлем если вы защищены от XSS.

### Q: Как защитить API токен?
**A:**
1. **Никогда** не коммитьте токены в git
2. **Никогда** не отправляйте токены в открытом виде
3. Используйте переменные окружения (`.env`)
4. Сохраняйте в защищенном хранилище
5. Регулярно ротируйте токены

Пример `.env`:
```
API_TOKEN=gog_1234567890abcdef
```

Использование:
```javascript
const token = process.env.API_TOKEN;
```

### Q: Что если мой API токен скомпрометирован?
**A:**
1. Немедленно отзовите токен через API
2. Создайте новый токен
3. Обновите токен во всех местах где он используется

---

## Интеграция

### Q: Как использовать API в React?
**A:** Используйте `fetch` или библиотеку `axios`:

```javascript
// fetch
const response = await fetch('http://localhost:8080/auth/login', {
  method: 'POST',
  headers: {'Content-Type': 'application/json'},
  body: JSON.stringify({email, password})
});

// или axios
import axios from 'axios';
const response = await axios.post('http://localhost:8080/auth/login', {
  email, password
});
```

### Q: Как использовать API в Vue?
**A:** Аналогично React, используйте `fetch` или `axios`:

```javascript
// Vue component
export default {
  methods: {
    async login(email, password) {
      const response = await fetch('http://localhost:8080/auth/login', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({email, password})
      });
      const data = await response.json();
      localStorage.setItem('jwtToken', data.user.token);
    }
  }
}
```

### Q: Как использовать API в Angular?
**A:** Используйте `HttpClient`:

```typescript
import { HttpClient } from '@angular/common/http';

@Injectable()
export class AuthService {
  constructor(private http: HttpClient) {}
  
  login(email: string, password: string) {
    return this.http.post('/auth/login', {email, password});
  }
}
```

### Q: Как использовать API в Node.js?
**A:** Используйте `node-fetch` или `axios`:

```javascript
const fetch = require('node-fetch');

async function login(email, password) {
  const response = await fetch('http://localhost:8080/auth/login', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({email, password})
  });
  return response.json();
}
```

---

## CORS

### Q: Получаю CORS ошибку при обращении к API
**A:** Это нормально для локальной разработки. CORS ограничения браузера не позволяют запросы с разных источников. Решения:

1. **Использовать Proxy** в development (в `package.json` для Create React App):
```json
{
  "proxy": "http://localhost:8080"
}
```

2. **Отключить CORS проверку** в браузере (только для разработки):
```bash
# Chrome на Windows
start chrome --disable-web-security --disable-gpu --user-data-dir=C:/ChromeTemp
```

3. **Включить CORS на сервере** (планируется):
```go
// Будет добавлено в будущих версиях
router.Use(cors.Default())
```

---

## Разработка

### Q: Как запустить сервер в режиме разработки?
**A:**
```bash
go run ./cmd/main/main.go
```

### Q: Как собрать production версию?
**A:**
```bash
go build -o build/server.exe ./cmd/main/main.go
```

### Q: Как разворачивать приложение?
**A:** Планируется добавить Docker поддержку:
```bash
docker build -t gogen-api .
docker run -p 8080:8080 gogen-api
```

### Q: Как сбросить БД?
**A:**
```bash
# Удалить все миграции (осторожно!)
goose -dir sql/migrations postgres "..." down

# Или переинициализировать контейнер
docker-compose down -v
docker-compose up
```

---

## Тестирование

### Q: Как протестировать API локально?
**A:**
1. Используйте Postman (импортируйте коллекцию из `docs/EXAMPLES.md`)
2. Используйте curl (примеры в `docs/API.md`)
3. Используйте REST Client расширение в VS Code

### Q: Как написать unit тесты?
**A:** Примеры будут добавлены в папку `tests/`

### Q: Как запустить интеграционные тесты?
**A:** Планируется добавить в будущих версиях.

---

## Поддержка

### Q: Где я могу получить помощь?
**A:**
- Проверьте эту документацию
- Посмотрите примеры в `docs/EXAMPLES.md`
- Откройте issue на GitHub
- Свяжитесь с поддержкой

### Q: Как сообщить об ошибке?
**A:** Откройте issue на GitHub с:
- Описанием проблемы
- Шагами воспроизведения
- Ожидаемым поведением
- Фактическим поведением
- Логами ошибок

### Q: Как предложить новую функцию?
**A:** Откройте issue с тегом `feature-request` и опишите:
- Что нужно добавить
- Почему это нужно
- Как это должно работать

---

## Версионирование

### Q: Какая версия API сейчас используется?
**A:** Версия 1.0 (стабильная)

### Q: Когда выйдет версия 2.0?
**A:** Планируется в Q3 2026 с новыми функциями.

### Q: Будут ли обратно совместимые изменения?
**A:** Да, версия 1.x будет получать обновления безопасности и исправления ошибок минимум 12 месяцев.
