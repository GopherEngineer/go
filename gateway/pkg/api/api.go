// API приложения "APIGateway".
package api

import (
	"bytes"
	"encoding/json"
	"gateway/pkg/logger"
	"gateway/pkg/types"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

const (
	// порт сервиса "Агрегатор Новостей"
	AGGREGATOR_PORT = "8081"
	// порт сервиса "Комментарии"
	COMMENTS_PORT = "8082"
	// порт сервиса "Модератор"
	MODERATOR_PORT = "8083"
)

type API struct {
	router *mux.Router
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
	// получение публикаций
	api.router.HandleFunc("/news", logger.Middleware(api.handlerNews)).Methods(http.MethodGet)
	// получение публикации по идентификатору
	api.router.HandleFunc("/news/{id}", logger.Middleware(api.handlerNewsDatailed)).Methods(http.MethodGet)
	// сохранение комментария для публикации
	api.router.HandleFunc("/comments/{id}", logger.Middleware(api.handlerCommentCreate)).Methods(http.MethodPost)
	// веб-приложение
	api.router.PathPrefix("/").Handler(http.StripPrefix("/", http.HandlerFunc(api.handler)))
}

// обертка над http.NewRequest для указания в запросе хедера
// с IP адресом клиента который делает запросы к APIGateway
func request(r *http.Request, method, url string, data []byte) (*http.Response, []byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}

	// подготовка HTTP заголовка запроса, что запрос будет типа JSON
	req.Header.Add("Content-Type", "application/json")
	// для передачи другим сервисам IP клиента
	req.Header.Add("X-Forwarded-For", r.RemoteAddr)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	data, err = io.ReadAll(res.Body)
	return res, data, err
}

// обработчик HTTP запросов для проксирования запроса к веб-приложению
func (api *API) handler(w http.ResponseWriter, r *http.Request) {
	// подготовка HTTP заголовка ответа сервера, что будет возвращен JSON
	w.Header().Set("Content-Type", "application/json")
	// подготовка HTTP заголовка ответа сервера, что можно делать запросы с любых хостов
	w.Header().Set("Access-Control-Allow-Origin", "*")

	res, err := http.Get("http://localhost:" + AGGREGATOR_PORT + "/" + r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
	w.Write(data)
}

// обработчик HTTP запросов получения списка публикаций
func (api *API) handlerNews(w http.ResponseWriter, r *http.Request) {
	res, data, err := request(r, "GET", "http://localhost:"+AGGREGATOR_PORT+"/"+r.URL.Path+"?"+r.URL.RawQuery, []byte{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
	w.Write(data)
}

// обработчик HTTP запросов получения публикации с комментариями
func (api *API) handlerNewsDatailed(w http.ResponseWriter, r *http.Request) {
	// получение идентификатора публикации
	id := mux.Vars(r)["id"]

	// подготовка переменной для формирования ответа на запрос
	post := types.NewsShortDetailed{}

	// канал для получения ответов на подзапросы
	ch := make(chan any, 2)

	// для приостановки
	wg := sync.WaitGroup{}
	wg.Add(2)

	// запускаем горутину для запроса публикации
	go func(ch chan<- any) {
		defer wg.Done()

		var post types.Post

		// запрос к сервису публикаций на получение одной публикации по ID публикации
		_, data, err := request(r, "GET", "http://localhost:"+AGGREGATOR_PORT+"/news/"+id+"?"+r.URL.RawQuery, nil)
		// в случае получения ошибки пишем её в канал
		if err != nil {
			ch <- err
			return
		}

		// раскодирование содержимого в структура типа Post
		json.Unmarshal(data, &post)
		// отправка полученной публикации в канал
		ch <- post
	}(ch)

	// запускаем горутину для запроса комментариев
	go func(ch chan<- any) {
		defer wg.Done()

		var comments []types.Comment

		// запрос к сервису комментриев на получение списка комментариев по ID публикации
		_, data, err := request(r, "GET", "http://localhost:"+COMMENTS_PORT+"/comments/"+id+"?"+r.URL.RawQuery, nil)
		// в случае получения ошибки пишем её в канал
		if err != nil {
			ch <- err
			return
		}

		// раскодирование содержимого в структура типа Comment
		json.Unmarshal(data, &comments)
		// отправка полученных комментариев в канал
		ch <- comments
	}(ch)

	// ждем пока запросы отрботают
	wg.Wait()
	close(ch)

	// перебираем в цикле ответы запросов к сервисам публикаций и комментариям
	// и с помощью приведение типов сохряняем результаты ответов от сервисов,
	// или в случае ошибок завершаем запрос с указанием ошибки
	for res := range ch {
		switch v := res.(type) {
		case types.Post:
			post.Post = v
		case []types.Comment:
			post.Comments = v
		case error:
			http.Error(w, v.Error(), http.StatusInternalServerError)
			return
		default:
			http.Error(w, "Неизвестная ошибка", http.StatusInternalServerError)
			return
		}
	}

	// указывам что ответ будет в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// кодируем структуру типа NewsShortDetailed
	data, err := json.Marshal(post)
	// обязательно обработка ошибки
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// пишем ответ
	w.Write(data)
}

// обработчик HTTP запросов сохранение комментария для публикации
func (api *API) handlerCommentCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	// обязательно обработка ошибки
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// запрос к сервису модерирования комментария
	res, _, err := request(r, "POST", "http://localhost:"+MODERATOR_PORT+"?"+r.URL.RawQuery, data)
	// обязательно обработка ошибки
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if res.StatusCode != 200 {
		// в случае если ошибки явно не было, но ответ не 200,
		// то завершам запрос с полученым кодом
		w.WriteHeader(res.StatusCode)
		w.Write(nil)
		return
	}

	// запрос к сервису комментариев для сохранения комментария
	res, _, err = request(r, "POST", "http://localhost:"+COMMENTS_PORT+r.URL.Path+"?"+r.URL.RawQuery, data)
	// обязательно обработка ошибки
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if res.StatusCode != 200 {
		// в случае если ошибки явно не было, но ответ не 200,
		// то завершам запрос с полученым кодом
		w.WriteHeader(res.StatusCode)
		w.Write(nil)
		return
	}

	// просто ответаем статусом 200
	w.Write(nil)
}
