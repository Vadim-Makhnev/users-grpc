# gRPC User Service

Простой gRPC сервис на Go для управления пользователями.

## Технологии
- Go
- gRPC
- Protocol Buffers
- Evans (gRPC client)
- PostgreSQL
- Docker + docker-compose

## Функциональность
- CreateUser - создание пользователя
- GetUser - получение пользователя по ID  
- ListUsers — список пользователей с пагинацией и сортировкой
- UpdateUser — частичное или полное обновление пользователя
- DeleteUser — удаление пользователя с возвратом данных  

## Генерация proto
```bash
make gen
```

## Запуск тестов 
```bash
make test
```

## Docker
```bash
# Запуск контейнера
make docker-run
# Остановка контейнера с очисткой
make docker-down
```

## Тестирование
```bash
# Подключение к серверу
evans -p 4000 proto/user.proto
```

# Тестирование с Evans
```bash
# Создать пользователя
call CreateUser
name => Vadim
email => vadim@example.com
age => 25

# Обновить только email (partial update)
call UpdateUser
id => 1
email::value => new@example.com

# Получить список (с пагинацией)
call ListUsers
page => 1
page_size => 10
sort => name
```

# Тестирование с grpcurl
```bash
# Создать
grpcurl -plaintext -d '{"name":"Test","email":"test@example.com","age":30}' \
  localhost:4000 user.UserService/CreateUser

# Обновить
grpcurl -plaintext -d '{"id":1,"email":"updated@example.com"}' \
  localhost:4000 user.UserService/UpdateUser
```
