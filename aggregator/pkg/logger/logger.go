// Пакет "Logger" для журналирования запросов.
package logger

import (
	"log/slog"
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

// Middleware оборачивает http обработчик
// для записи в лог информации о запросе.
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// испольуем обертку над http.ResponseWriter
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)
		// пишем в лог информацию о запросе
		slog.Info(r.URL.RequestURI(), "Метод", r.Method, "IP-адрес", r.Header.Get("X-Forwarded-For"), "HTTP-код ответа", rw.statusCode, "ID запроса", r.URL.Query().Get("request_id"))
	}
}
