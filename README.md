# Currency & Weather Aggregator API (Go)

Сервис на Go, агрегирующий данные о погоде и курсах валют с внешних API и предоставляющий REST API.

## Эндпоинты
- GET `/weather?city=City`
- GET `/currency?base=USD&target=EUR`
- GET `/report?city=City&base=USD&target=EUR`

## Стек и архитектура
- Go + Gin, Clean Architecture (handler → service → provider), DTO в `internal/model`
- Параллельные запросы (goroutines + channels), таймауты (context), retries (exponential backoff)
- Кэш in-memory с TTL, логирование через zap, конфигурация через viper

## Конфигурация (.env)
```
PORT=8081
WEATHER_API_KEY=...   # ваш ключ OpenWeather
WEATHER_BASE_URL=https://api.openweathermap.org/data/2.5
RATES_BASE_URL=https://open.er-api.com/v6
REQUEST_TIMEOUT=5s
CACHE_TTL=30s
LOG_LEVEL=info
```

## Запуск
```
go run ./cmd
```

## Примеры запросов
```
curl "http://localhost:8081/weather?city=London"
curl "http://localhost:8081/currency?base=USD&target=EUR"
curl "http://localhost:8081/report?city=London&base=USD&target=EUR"
```

## Тесты
```
go test ./...
```

## Дальнейшее развитие
- Swagger/OpenAPI, Prometheus метрики, Docker Compose, CI расширение
