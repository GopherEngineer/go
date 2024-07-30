// API сервера "Комментарии".
package api

import (
	"comments/pkg/logger"
	"comments/pkg/storage"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type API struct {
	db     storage.Interface
	router *mux.Router
}

// Конструктор API.
func New(db storage.Interface) *API {
	api := API{
		db:     db,
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
	// получение комментариев для публикации
	api.router.HandleFunc("/comments/{id}", logger.Middleware(api.handler)).Methods(http.MethodGet)
	// сохранение комментария для публикации
	api.router.HandleFunc("/comments/{id}", logger.Middleware(api.handlerCreate)).Methods(http.MethodPost)
}

// обработчик HTTP запросов получения комментариев для публикации
func (api *API) handler(w http.ResponseWriter, r *http.Request) {
	// получение идентификатора публикации
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	// обязательно обработка ошибки
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// получение комментариев из БД по идентификатору публикации
	comments, err := api.db.Comments(id)

	// обязательно обработка ошибки
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// пишем ответ кодируя комментарии
	json.NewEncoder(w).Encode(comments)

}

// обработчик HTTP запросов сохранение комментария для публикации
func (api *API) handlerCreate(w http.ResponseWriter, r *http.Request) {
	// получение идентификатора публикации
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	// обязательно обработка ошибки
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	// чтение содержимого полученного запроса
	data, err := io.ReadAll(r.Body)
	// обязательно обработка ошибки
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// раскодирование содержимого в структура типа Comment
	var comment storage.Comment
	err = json.Unmarshal(data, &comment)
	// обязательно обработка ошибки
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// указание ID публикации непосредственно из полученного идентификатора
	comment.Post = id
	// указание текущей даты
	comment.CreatedAt = time.Now().Unix()

	// сохранение комметария в БД
	err = api.db.SaveComment(comment)
	// обязательно обработка ошибки
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// просто ответаем статусом 200
	w.Write(nil)
}
