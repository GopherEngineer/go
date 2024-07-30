// API приложения "Новостной агрегатор".
package api

import (
	"aggregator/pkg/logger"
	"aggregator/pkg/storage"
	"encoding/json"
	"math"
	"net/http"
	"strconv"

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
	// получение публикаций
	api.router.HandleFunc("/news", logger.Middleware(api.handler)).Methods(http.MethodGet)
	// получение публикации по ID
	api.router.HandleFunc("/news/{id}", logger.Middleware(api.handlerDetail)).Methods(http.MethodGet)
	// веб-приложение
	api.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./web/public"))))
}

// обработчик HTTP запросов получения публикаций
func (api *API) handler(w http.ResponseWriter, r *http.Request) {
	// подготовка HTTP заголовка ответа сервера, что будет возвращен JSON
	w.Header().Set("Content-Type", "application/json")
	// подготовка HTTP заголовка ответа сервера, что можно делать запросы с любых хостов
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// получение из запроса номера страницы публикаций
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	// получение параметра для поиска
	search := r.URL.Query().Get("s")

	// обращние к базе данных для получения публикаций
	news, n, err := api.db.News(page, search)
	// если произошла ошибка чтения из базы данных, то завершаем
	// обращение к серверу и указываем, что произошла
	// внутренняя ошибка сервера с текстом ошибки
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nfd := storage.NewsFullDetailed{
		News:  news,
		Page:  page,
		Pages: int(math.Ceil(float64(n) / 15.0)),
		Count: len(news),
	}

	// пишем ответ кодируя публикации
	json.NewEncoder(w).Encode(nfd)
}

// обработчик HTTP запросов получения публикации по ID
func (api *API) handlerDetail(w http.ResponseWriter, r *http.Request) {
	// подготовка HTTP заголовка ответа сервера, что будет возвращен JSON
	w.Header().Set("Content-Type", "application/json")
	// подготовка HTTP заголовка ответа сервера, что можно делать запросы с любых хостов
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// получение идентификатора публикации
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// обращние к базе данных для получения публикации
	post, err := api.db.NewsDetail(id)
	// если произошла ошибка чтения из базы данных, то завершаем
	// обращение к серверу и указываем, что произошла
	// внутренняя ошибка сервера с текстом ошибки
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// если публикации нет, то возвращаем пустой объект в JSON спецификации
	if post == nil {
		w.Write([]byte("{}"))
		return
	}

	// пишем ответ кодируя публикацию
	json.NewEncoder(w).Encode(post)
}
