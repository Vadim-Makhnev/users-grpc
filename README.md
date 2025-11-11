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

## Запуск

Генерация кода из .proto:
```bash
make gen