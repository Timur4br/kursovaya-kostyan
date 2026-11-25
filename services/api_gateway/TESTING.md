# Проверка API

Из корня репозитория, при запущенном Docker:

```sh
docker compose -f docker-compose.test.yaml up --abort-on-container-exit --exit-code-from tests
```

Успешный запуск завершается с кодом 0. После проверки:

```sh
docker compose -f docker-compose.test.yaml down
```

Проверка использует отдельный PostgreSQL с pgvector, без публикации портов и без
постоянного хранения данных. Настройки рабочей базы из `.env` не используются.
Первый запуск скачивает образы и зависимости Go.

Тесты проверяют маршруты `camera`, `object` и `event` через controller → service → pg_dao:
создание, получение, полное обновление, удаление, списки, пагинацию, количество,
отсутствующие записи, внешние ключи и передачу вектора `vector(512)`.
Схема временных таблиц берётся из `scripts/init-db.sql`.

Для локального запуска Go-проверок из `services/api_gateway/src`:

```sh
go test ./...
go vet ./...
```

Без `TEST_DATABASE_URL` тест PostgreSQL пропускается. Для его запуска нужна
строка подключения к PostgreSQL с установленным расширением `vector`.
Тест создаёт только временные таблицы на одном соединении; при закрытии
соединения они удаляются.

## Новые маршруты

Для `object` и `event`, как для `camera`, доступны:

```text
POST   /api/object
GET    /api/object/:id
PUT    /api/object/:id
DELETE /api/object/:id
GET    /api/object/list?limit=10&offset=0
GET    /api/object/count

POST   /api/event
GET    /api/event/:uuid
PUT    /api/event/:uuid
DELETE /api/event/:uuid
GET    /api/event/list?limit=10&offset=0
GET    /api/event/count
```

`objects` хранит уникальный `type`. `events` содержит `camera_id`, `object_id`,
`type`, `confidence`, `embedding` (массив из 512 чисел) и `photo_uuid`.
При создании идентификатор и время назначает PostgreSQL. `PUT` заменяет все
изменяемые поля; идентификатор берётся из URL, время создания сохраняется.
Полный контракт доступен в Swagger по `/swagger/index.html`.

## Опциональные фильтры

Фильтры передаются query-параметрами в `/list` и `/count`. У обоих методов
одинаковые условия; `limit` и `offset` влияют только на список.
Все заданные условия объединяются через AND. При отсутствии фильтров
возвращаются все записи с учётом пагинации.

| Сущность | Параметры |
|---|---|
| camera | `camera_id`, `ip`, `is_enabled` |
| object | `object_id`, `type` |
| event | `camera_id`, `object_id`, `type`, `confidence_min`, `confidence_max` |
| Все три | `created_from`, `created_to` |

`type` и `ip` сравниваются точно. Границы confidence и времени включительные.
Для времени подходят RFC3339 (с часовым поясом) и форматы из `setDateParam`;
дата без времени означает 00:00:00 UTC, а не весь день.
Например, часовой пояс `+03:00` в URL следует записывать как `%2B03:00`.
Неверные типы значений и обратные границы диапазонов дают HTTP 400.

```http
GET /api/event/list?camera_id=1&confidence_min=0.8&limit=10&offset=0
GET /api/event/count?camera_id=1&confidence_min=0.8
GET /api/camera/list?is_enabled=false
GET /api/object/list?type=backpack
```

## Изменение и удаление камер

`PUT /api/camera/:id` принимает `entities.CameraUpdate`. Пропущенные поля
не меняются, пустые строки очищают название/адрес потока, `is_enabled: false`
отключает камеру. ID из URL и время создания сохраняются.
Порт должен быть от 1 до 65535; IP/хост не может быть пустым.

`DELETE /api/camera/:id` удаляет камеру без событий. При наличии событий
возвращается 409, запись и история сохраняются. Для отсутствующей камеры
GET, PUT и DELETE возвращают 404.
