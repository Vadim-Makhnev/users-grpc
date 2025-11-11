# gRPC User Service

Простой gRPC сервис на Go для управления пользователями.

## Технологии
- Go
- gRPC
- Protocol Buffers
- Evans (gRPC client)

## Функциональность
- CreateUser - создание пользователя
- GetUser - получение пользователя по ID  
- ListUsers - список всех пользователей

## Генерация proto
```bash
make gen
```

## Тестирование с Evans
```bash
# Подключение к серверу
evans -p 50051 proto/user.proto
```

### Примеры запросов:

```bash
# Создание пользователя
call CreateUser
name (TYPE_STRING) => Vadim
email (TYPE_STRING) => vadim@example.com  
age (TYPE_INT32) => 21

# Ответ:
{
  "age": 21,
  "email": "vadim@example.com",
  "id": "ac8b2c6c-1969-4c57-9d93-400cf7ef05c4",
  "name": "Vadim"
}
```

```bash
# Получение пользователя по ID
call GetUser
id (TYPE_STRING) => ac8b2c6c-1969-4c57-9d93-400cf7ef05c4

# Ответ:
{
  "age": 21,
  "email": "vadim@example.com", 
  "id": "ac8b2c6c-1969-4c57-9d93-400cf7ef05c4",
  "name": "Vadim"
}
```

```bash  
# Получение списка всех пользователей
call ListUsers

# Ответ:
{
  "users": [
    {
      "age": 21,
      "email": "vadim@example.com",
      "id": "ac8b2c6c-1969-4c57-9d93-400cf7ef05c4",
      "name": "Vadim"
    }
  ]
}
```