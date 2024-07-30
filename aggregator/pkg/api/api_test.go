package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"aggregator/pkg/storage"
	"aggregator/pkg/storage/memdb"
)

func request(api *API, method, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, nil)

	rr := httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)

	return rr
}

func Test(t *testing.T) {

	db, _ := memdb.New()
	api := New(db)

	t.Run("news", func(t *testing.T) {
		rr := request(api, "GET", "/news")

		if http.StatusOK != rr.Code {
			t.Errorf("Ожидалался код ответа %d. Получено %d\n", http.StatusOK, rr.Code)
		}

		body := rr.Body.Bytes()

		var full storage.NewsFullDetailed

		json.Unmarshal(body, &full)

		if len(full.News) != 15 {
			t.Errorf("Ожидалось получить 15 публикаций. Получили %d", len(full.News))
		}
	})

	t.Run("news/1", func(t *testing.T) {
		rr := request(api, "GET", "/news/1")

		if http.StatusOK != rr.Code {
			t.Errorf("Ожидалался код ответа %d. Получено %d\n", http.StatusOK, rr.Code)
		}

		body := rr.Body.Bytes()

		var post *storage.Post

		json.Unmarshal(body, &post)

		if post == nil || post.ID != 1 {
			t.Error("Ожидалось получить 1 публикацию. Получили 0")
		}
	})
}
