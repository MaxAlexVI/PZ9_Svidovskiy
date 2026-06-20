# Практическое задание №9

## Тема

Redis-кэш для учебного REST-сервиса `tasks`.

## Цель работы

Добавить к REST-сервису слой кэширования через Redis, ускорить повторное чтение задач и реализовать инвалидацию кэша при изменении данных.

## Структура проекта

```
pz9-redis-cache/
├── cmd/
│   └── server/
│       └── main.go
├── deploy/
│   └── redis/
│       └── docker-compose.yml
├── internal/
│   ├── cache/
│   ├── config/
│   ├── httpapi/
│   ├── service/
│   └── task/
├── go.mod
├── go.sum
└── README.md
```

## Как работает

При запросе задачи сервис сначала проверяет Redis. Если данные найдены, ответ формируется из кэша. Если данных нет, сервис читает задачу из in-memory репозитория, возвращает ее клиенту и сохраняет в Redis с ограниченным временем жизни. При изменении или удалении задачи сервис очищает связанные ключи, чтобы клиент не получил устаревшие данные.

## Маршруты

```
GET    /v1/tasks
POST   /v1/tasks
GET    /v1/tasks/{id}
PATCH  /v1/tasks/{id}
DELETE /v1/tasks/{id}
```

## Запуск Redis

```
cd deploy\redis
docker compose up -d
```

<img width="823" height="139" alt="image" src="https://github.com/user-attachments/assets/6d19ee73-6a2e-4078-8226-337d81162ca5" />

## Запуск сервиса

```
go run .\cmd\server
```

<img width="508" height="38" alt="image" src="https://github.com/user-attachments/assets/76797d06-d80d-4c55-8adb-be698f725155" />

## Проверка
`GET`
<img width="1698" height="218" alt="image" src="https://github.com/user-attachments/assets/44249e58-c319-48a2-b5af-88f1cc8a7dcc" />
<img width="1697" height="242" alt="image" src="https://github.com/user-attachments/assets/2528f02b-40b9-4fbd-adfc-151212d2000c" />

`POST`
<img width="1698" height="293" alt="image" src="https://github.com/user-attachments/assets/682f0d60-f243-4f75-bf95-315e24aa633c" />

`PATCH`
<img width="1700" height="335" alt="image" src="https://github.com/user-attachments/assets/92dc28ae-10a2-4b67-8e3c-14988681ee62" />

`DELETE`
<img width="1698" height="291" alt="image" src="https://github.com/user-attachments/assets/1bc1e876-78dc-41f4-9293-fa24fcd1a569" />

## Ожидаемый результат

Первый `GET /v1/tasks/1` читает данные из репозитория и кладет их в Redis. Повторный запрос того же id может обслуживаться из кэша. После `PATCH` или `DELETE` связанные кэш-ключи удаляются.

## Вывод

В ходе работы REST-сервис был дополнен Redis-кэшем. Практика показывает типовой подход к ускорению чтения и поддержанию актуальности данных через инвалидацию.
