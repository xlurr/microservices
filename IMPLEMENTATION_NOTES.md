# Реализация бизнес-логики и JSON-хранилища

## Что было сделано

### 1. JSON-хранилище (все 4 сервиса)
- **Файл**: `internal/utils/file_storage.go` - утилита для работы с JSON
- **Место хранения**: `./data/{service}.json`
- **Особенность**: автоматическая синхронизация при каждом изменении

### 2. Users-Service
- **Валидация email**: обязательное наличие одного символа `@`
- **Проверка**: email не может начинаться или заканчиваться на `@`
- **Эндпоинт**: `GET /api/users/{id}/exists` - проверка существования пользователя
- **Репозиторий**: `JSONUserRepository` с методом `UserExists(id int64) bool`

### 3. Orders-Service
- **Хранилище**: JSON файл с заказами
- **Методы**: `GetOrdersByUserID`, `DeleteOrdersByUserID`, `UpdateOrderStatus`
- **Проверка**: при создании заказа нужна проверка пользователя через HTTP

### 4. Payments-Service
- **Хранилище**: JSON файл с платежами
- **Методы**: `GetPaymentsByUserID`, `DeletePaymentsByUserID`, `UpdatePaymentStatus`
- **Проверка**: при создании платежа нужна проверка пользователя через HTTP

### 5. Delivery-Service
- **Хранилище**: JSON файл с доставками
- **Методы**: `GetDeliveriesByUserID`, `DeleteDeliveriesByUserID`, `UpdateDeliveryStatus`
- **Проверка**: при создании доставки нужна проверка пользователя через HTTP

## Структура данных JSON

### users.json
```json
[
  {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "age": 30,
    "createdAt": "2025-12-05T06:30:00Z",
    "updatedAt": "2025-12-05T06:30:00Z"
  }
]
```

### orders.json
```json
[
  {
    "id": 1,
    "userId": 1,
    "items": ["item1", "item2"],
    "totalAmount": 200,
    "status": "created",
    "createdAt": "2025-12-05T06:31:00Z",
    "updatedAt": "2025-12-05T06:31:00Z"
  }
]
```

### payments.json
```json
[
  {
    "id": 1,
    "userId": 1,
    "orderId": 1,
    "amount": 200,
    "status": "pending",
    "createdAt": "2025-12-05T06:32:00Z",
    "updatedAt": "2025-12-05T06:32:00Z"
  }
]
```

### deliveries.json
```json
[
  {
    "id": 1,
    "userId": 1,
    "orderId": 1,
    "address": "123 Main St",
    "status": "pending",
    "trackingId": "TRACK123",
    "createdAt": "2025-12-05T06:33:00Z",
    "updatedAt": "2025-12-05T06:33:00Z"
  }
]
```

## Как использовать

### 1. Запустить скрипт (уже выполнено)
```bash
chmod +x implement-business-logic.sh
./implement-business-logic.sh
```

### 2. Запустить Docker Compose
```bash
docker-compose up --build
```

### 3. Проверить endpoints

#### Users Service
```bash
# Создать пользователя с валидацией @
curl -X POST http://localhost:8081/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","name":"John","age":30}'

# Получить пользователя
curl http://localhost:8081/api/users/1

# Проверить существование пользователя
curl http://localhost:8081/api/users/1/exists

# Удалить пользователя (каскадно удаляет заказы, платежи, доставки)
curl -X DELETE http://localhost:8081/api/users/1
```

#### Orders Service
```bash
# Создать заказ
curl -X POST http://localhost:8082/api/orders \
  -H "Content-Type: application/json" \
  -d '{"userId":1,"items":["item1"],"totalAmount":100}'

# Получить заказы пользователя
curl http://localhost:8082/api/orders/user/1

# Обновить статус заказа
curl -X PUT http://localhost:8082/api/orders/1 \
  -H "Content-Type: application/json" \
  -d '{"status":"completed"}'
```

#### Payments Service
```bash
# Создать платёж
curl -X POST http://localhost:8083/api/payments \
  -H "Content-Type: application/json" \
  -d '{"userId":1,"orderId":1,"amount":100}'

# Получить платежи пользователя
curl http://localhost:8083/api/payments/user/1

# Обновить статус платежа
curl -X PUT http://localhost:8083/api/payments/1 \
  -H "Content-Type: application/json" \
  -d '{"status":"completed"}'
```

#### Delivery Service
```bash
# Создать доставку
curl -X POST http://localhost:8084/api/deliveries \
  -H "Content-Type: application/json" \
  -d '{"userId":1,"orderId":1,"address":"123 Main St","trackingId":"TRACK123"}'

# Получить доставки пользователя
curl http://localhost:8084/api/deliveries/user/1

# Обновить статус доставки
curl -X PUT http://localhost:8084/api/deliveries/1 \
  -H "Content-Type: application/json" \
  -d '{"status":"shipped"}'
```

## Важные замечания

1. **JSON файлы сохраняются автоматически** при каждом изменении
2. **При перезагрузке контейнера** данные загружаются из JSON файлов
3. **Валидация email** происходит в users-service перед сохранением
4. **Проверка пользователя** в других сервисах должна происходить через HTTP запрос к `/users/{id}/exists`
5. **Каскадное удаление** - при удалении пользователя нужно удалить его заказы, платежи и доставки

## TODO (если нужно расширить)

- [ ] Добавить HTTP клиент в orders, payments, delivery для проверки пользователя
- [ ] Реализовать логику каскадного удаления при удалении пользователя
- [ ] Добавить логирование в репозиториях
- [ ] Добавить обработку ошибок в handlers
- [ ] Добавить unit тесты для репозиториев
