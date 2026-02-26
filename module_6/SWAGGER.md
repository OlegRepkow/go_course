# Swagger Documentation

## Генерація документації

```bash
make swagger
# або
~/go/bin/swag init -g cmd/server/main.go
```

## Доступ до Swagger UI

Після запуску сервера відкрийте:
```
http://localhost:8080/swagger/index.html
```

## Використання

1. Спочатку виконайте `POST /auth/sign-up` для реєстрації
2. Потім `POST /auth/sign-in` для отримання токена
3. Натисніть "Authorize" у Swagger UI та введіть: `Bearer YOUR_TOKEN`
4. Тепер можете використовувати захищені endpoints

## Endpoints

- **POST /auth/sign-up** - Реєстрація нового користувача
- **POST /auth/sign-in** - Вхід та отримання JWT токена
- **GET /channel/history** - Отримання історії повідомлень (потрібна авторизація)
- **POST /channel/send** - Відправка повідомлення (потрібна авторизація)

## Оновлення документації

Після змін в коді handlers:
```bash
make swagger
```

Документація автоматично оновиться в папці `docs/`.
