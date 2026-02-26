# Chat Service

Мінімальний чат-сервіс з REST API та WebSocket.

## Запуск

```bash
# Завантажити залежності
go mod tidy

# Згенерувати Swagger документацію
make swagger

# Запустити через Docker
docker-compose up --build

# Або локально
go run cmd/server/main.go
```

## API Documentation

Swagger UI доступний за адресою: http://localhost:8080/swagger/index.html

## Endpoints

- `POST /auth/sign-up` - реєстрація (body: `{"username":"user","password":"pass"}`)
- `POST /auth/sign-in` - вхід (body: `{"username":"user","password":"pass"}`, повертає token)
- `GET /channel/history` - історія повідомлень (потрібен Bearer token)
- `POST /channel/send` - надіслати повідомлення (body: `{"text":"message"}`, потрібен Bearer token)
- `WS /channel/listen` - WebSocket для отримання повідомлень в реальному часі

## Приклад використання

```bash
# Реєстрація
curl -X POST http://localhost:8080/auth/sign-up \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'

# Вхід
curl -X POST http://localhost:8080/auth/sign-in \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'

# Надіслати повідомлення
curl -X POST http://localhost:8080/channel/send \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"text":"Hello!"}'

# Історія
curl http://localhost:8080/channel/history \
  -H "Authorization: Bearer YOUR_TOKEN"
```
