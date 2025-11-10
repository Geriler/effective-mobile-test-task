# Задача

Спроектировать и реализовать REST-сервис для агрегации данных об онлайн-подписках пользователей.

## Требования

1. Выставить HTTP-ручки для CRUDL-операций над записями о подписках. Каждая запись содержит:
   1. Название сервиса, предоставляющего подписку
   2. Стоимость месячной подписки в рублях
   3. ID пользователя в формате UUID
   4. Дата начала подписки (месяц и год)
   5. Опционально дата окончания подписки
2. Выставить HTTP-ручку для подсчета суммарной стоимости всех подписок за выбранный период с фильтрацией по id пользователя и названию подписки
3. СУБД – PostgreSQL. Должны быть миграции для инициализации базы данных
4. Покрыть код логами
5. Вынести конфигурационные данные в .env/.yaml-файл
6. Предоставить swagger-документацию к реализованному API
7. Запуск сервиса с помощью docker compose

## Примечания:
1. Проверка существования пользователя не требуется. Управление пользователями вне зоны ответственности вашего сервиса
2. Стоимость любой подписки – целое число рублей, копейки не учитываются

Пример тела запроса на создание записи о подписке:
```json
{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
}
```

# Быстрый старт

## Запуск

```bash
docker compose up
```

Команда автоматически:

- Разворачивает PostgreSQL
- Применяет миграции
- Запускает gRPC сервер на порту **8081**
- Запускает HTTP Gateway на порту **8080**

### Проверка работоспособности

После запуска сервис будет доступен по адресу: [http://localhost:8080](http://localhost:8080)

### Swagger UI

Интерактивная документация API: [http://localhost:8080/swagger](http://localhost:8080/swagger)

### Пример запроса (curl)

```bash
# Создание подписки
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'

# Получение списка подписок
curl http://localhost:8080/api/v1/subscriptions

# Подсчет суммы за период
curl "http://localhost:8080/api/v1/subscriptions/sum?startDate=01-2025&endDate=12-2025"
```

Больше примеров: [example/subscriptions.http](example/subscriptions.http)

# API Endpoints

## Подписки

- `POST /api/v1/subscriptions` - Создать подписку

- `GET /api/v1/subscriptions` - Получить все подписки

- `GET /api/v1/subscriptions/{id}` - Получить подписку по ID

- `PUT /api/v1/subscriptions/{id}` - Обновить подписку

- `DELETE /api/v1/subscriptions/{id}` - Удалить подписку

## Аналитика

- `GET /api/v1/subscriptions/sum` — Подсчет суммы с фильтрацией

**Параметры фильтрации:**

- `startDate` - дата начала периода (формат: `MM-YYYY`)
- `endDate` - дата окончания периода (формат: `MM-YYYY`)
- `serviceName` - наименование подписки (опционально)
- `userId` - ID пользователя (опционально)

# Конфигурация

Файлы конфигураций находятся в [configs/config.yml](configs/config.yml) и [configs/.env](configs/.env)

# Миграции

Миграции применяются автоматически при запуске через Docker Compose

Для ручного применения:

```bash
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/postgres" up
```

# Генерация кода

```bash
# Protobuf + gRPC + Gateway + Swagger
buf dep update && buf generate

# SQLC
sqlc generate
```

# Локальный запуск (без Docker compose)

```bash
# Запустить PostgreSQL
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres:18

# Применить миграции
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/postgres" up

# Запустить сервер
go run cmd/server/main.go --config=configs/local.yml
```