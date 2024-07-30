// API сервера "Модератор комментариев".
package api

import (
	"encoding/json"
	"io"
	"moderator/pkg/logger"
	"moderator/pkg/moderator"
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	router *mux.Router
}

type Comment struct {
	Message string `json:"message"`
}

// Конструктор API.
func New() *API {
	api := API{
		router: mux.NewRouter(),
	}
	api.endpoints()
	return &api
}

// Router возвращает маршрутизатор для использования
// в качестве аргумента HTTP-сервера.
func (api *API) Router() *mux.Router {
	return api.router
}

// регистрация методов API в маршрутизаторе запросов
func (api *API) endpoints() {
	// модерирование комментария
	api.router.HandleFunc("/", logger.Middleware(api.handler)).Methods(http.MethodPost)
}

// обработчик HTTP запросов модерирование комментария
func (api *API) handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// чтение содержимого полученного запроса
	data, err := io.ReadAll(r.Body)
	// обязательно обработка ошибки
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// раскодирование содержимого в структура типа Comment
	var body Comment
	err = json.Unmarshal(data, &body)
	// обязательно обработка ошибки
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// проверка текста комментария на запрещенные слова
	if !moderator.Moderate(body.Message) {
		// если проверка не прошла, что отвечаем статусом 400
		w.WriteHeader(http.StatusBadRequest)
		w.Write(nil)
		return
	}

	// просто ответаем статусом 200
	w.Write(nil)
}
