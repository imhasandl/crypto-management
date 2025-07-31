# Crypto Management Service

## Запуск проекта

1. Установите Docker и Docker Compose.
2. Склонируйте репозиторий и перейдите в папку проекта.
3. Запустите сервисы:
   ```sh
   docker compose up --build

   или

   docker-compose up --build
   ```

## Установка зависимостей

Все зависимости устанавливаются автоматически при сборке Docker-образа (`go mod download`).

## Поддерживаемые монеты

В запросах используйте идентификаторы монет CoinGecko. Примеры популярных монет:
- `bitcoin`
- `ethereum`
- `dogecoin`
- `solana`
- `cardano`
- `litecoin`

Полный список доступен на [CoinGecko API](https://api.coingecko.com/api/v3/coins/list).

## Примеры curl-запросов

### Добавить отслеживание монеты

> Эндпоинт `/currency/add` начинает отслеживать монету и получает её цену каждые 10 секунд.

```sh
curl -X POST http://localhost:8080/currency/add \
  -H "Content-Type: application/json" \
  -d '{"coin":"bitcoin"}'
```

### Удалить монету из отслеживания

> Эндпоинт `/currency/remove` останавливает runner для указанной монеты и прекращает получение её данных.

```sh
curl -X POST http://localhost:8080/currency/remove \
  -H "Content-Type: application/json" \
  -d '{"coin":"bitcoin"}'
```

### Получить цену монеты на определённое время

> Эндпоинт `/currency/price` возвращает цену монеты, ближайшую к указанному времени (timestamp).  
> В теле запроса указывается идентификатор монеты и Unix-время.  
> В ответе приходит цена и фактическое время, к которому она относится.

```sh
curl -X POST http://localhost:8080/currency/price \
  -H "Content-Type: application/json" \
  -d '{"coin":"bitcoin", "timestamp": 1718000000}'
```

---
