# libs

Набор небольших Go-компонентов для HTTP-сервисов:

- `http_server` — обёртка над `net/http` с конфигурацией сервера, graceful shutdown, группами маршрутов и middleware;
- `rate_limiter/sliding_window` — ограничение количества запросов в скользящем временном окне;
- `rate_limiter/token_bucket` — token bucket с постоянным пополнением токенов.

Модуль использует Go 1.25.5.

## Установка

```bash
go get github.com/n1k1x86/libs
```

## HTTP-сервер

### Минимальный пример

```go
package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	httpserver "github.com/n1k1x86/libs/http_server"
)

func main() {
	mux := httpserver.NewMux()
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("hello world"))
	})

	cfg := httpserver.NewHTTPServerConfig().
		WithAddr(":8080").
		WithReadHeaderTimeout(5 * time.Second).
		WithReadTimeout(10 * time.Second).
		WithWriteTimeout(10 * time.Second).
		WithIdleTimeout(60 * time.Second)

	server := httpserver.NewHTTPServer(cfg).WithMux(mux)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Start()
	}()

	select {
	case err := <-errCh:
		if err != nil {
			log.Fatal(err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
	}
}
```

`Start` вызывает `http.Server.ListenAndServe`. Штатное завершение через `Shutdown` не возвращает `http.ErrServerClosed` как ошибку.

### Конфигурация

Конфигурация строится цепочкой методов, каждый из которых возвращает изменённую копию:

```go
cfg := httpserver.NewHTTPServerConfig().
	WithAddr(":8080").
	WithDisableGeneralOptionsHandler(false).
	WithReadHeaderTimeout(5 * time.Second).
	WithReadTimeout(10 * time.Second).
	WithWriteTimeout(10 * time.Second).
	WithIdleTimeout(60 * time.Second)
```

Неуказанные поля сохраняют нулевые значения `net/http.Server`.

## Middleware

Middleware имеет стандартную для `net/http` сигнатуру:

```go
type Middleware func(next http.Handler) http.Handler
```

Пример middleware:

```go
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
```

Функция `Chain` собирает цепочку вокруг конечного обработчика:

```go
handler := httpserver.Chain(finalHandler, logging, authentication)
```

Middleware выполняются в порядке передачи: сначала `logging`, затем `authentication`, затем `finalHandler`. После возврата из обработчика управление идёт в обратном порядке.

## Группы маршрутов

Группа добавляет общий префикс к маршрутам и применяет к ним общий набор middleware:

```go
api := mux.Group("/api", logging, authentication)

api.HandleFunc("GET /users", getUsers)
api.HandleFunc("POST /users", createUser)
api.Handle("DELETE /users/{id}", deleteUserHandler)
```

В результате регистрируются маршруты:

```text
GET /api/users
POST /api/users
DELETE /api/users/{id}
```

Префикс следует передавать без завершающего `/`, а паттерн маршрута — с `/` в начале пути.

Паттерны маршрутов передаются стандартному `http.ServeMux`, поэтому доступны HTTP-методы и path wildcards из актуального `net/http`:

```go
api.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_, _ = w.Write([]byte(id))
})
```

## Sliding window rate limiter

Пакет ограничивает число разрешённых операций для каждого `userID` за заданное временное окно.

```go
import (
	"time"

	"github.com/n1k1x86/libs/rate_limiter/sliding_window"
)

limiter := sliding_window.New(100, time.Minute)

if !limiter.Allow(userID) {
	// Лимит в 100 операций за последнюю минуту исчерпан.
}
```

Параметры `New`:

- `limit` — максимальное количество разрешённых операций в окне;
- `window` — длительность скользящего окна.

Состояние хранится отдельно для каждого `userID`. Вызов `Allow` потокобезопасен.

## Token bucket rate limiter

Token bucket допускает кратковременные всплески нагрузки до размера bucket и затем восстанавливает токены с заданной скоростью.

```go
import "github.com/n1k1x86/libs/rate_limiter/token_bucket"

// 5 токенов в секунду, максимальный запас — 20 токенов.
limiter := token_bucket.New(5, 20)

if !limiter.Allow(userID) {
	// В bucket пользователя сейчас нет доступного токена.
}
```

Параметры `New`:

- `rate` — количество восстанавливаемых токенов в секунду;
- `capacity` — максимальное количество токенов в bucket.

Новый пользователь получает полностью заполненный bucket. Каждая успешная операция расходует один токен. Состояние хранится отдельно для каждого `userID`, а `Allow` потокобезопасен.

## Проверка

Запуск тестов и проверки сборки всех пакетов:

```bash
go test ./...
```

