// Пакет "Logger" для журналирования запросов.
package logger

import (
	"log/slog"
	"math/rand"
	"net/http"
)

// обертка над http.ResponseWriter для
// извлечения статуса ответа
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// для генерации ID строкового типа с проивольной длиной
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Middleware оборачивает http обработчик
// для записи в лог информации о запросе.
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// подготовка HTTP заголовка ответа сервера, что будет возвращен JSON
		w.Header().Set("Content-Type", "application/json")
		// подготовка HTTP заголовка ответа сервера, что можно делать запросы с любых хостов
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// если ID запроса не был передан извне, то добавляем его
		if r.URL.Query().Get("request_id") == "" {
			q := r.URL.Query()
			q.Set("request_id", randString(6))
			r.URL.RawQuery = q.Encode()
		}

		// испольуем обертку над http.ResponseWriter
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)

		// пишем в лог информацию о запросе
		slog.Info(r.URL.RequestURI(), "Метод", r.Method, "IP-адрес", r.RemoteAddr, "HTTP-код ответа", rw.statusCode, "ID запроса", r.URL.Query().Get("request_id"))

	}
}
