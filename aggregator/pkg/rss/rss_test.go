package rss

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestParse(t *testing.T) {
	news, err := Parse("https://habr.com/ru/rss/hubs/go/articles/?fl=ru")
	if err != nil {
		t.Fatal(err)
	}
	if len(news) == 0 {
		t.Fatal("Количество публикация после раскодирования должно быть больше 0")
	}
}

func TestMockParse(t *testing.T) {
	rss, err := os.ReadFile("../../rss.xml")
	if err != nil {
		t.Fatal(err)
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, string(rss))
	}
	r := mux.NewRouter()
	r.HandleFunc("/", h).Methods(http.MethodGet)
	server := http.Server{
		Addr:    ":5000",
		Handler: r,
	}

	t.Cleanup(func() {
		server.Shutdown(context.Background())
	})
	go server.ListenAndServe()

	time.Sleep(time.Millisecond * 100)

	news, err := Parse("http://localhost:5000")
	if err != nil {
		t.Fatal(err)
	}
	if len(news) == 0 {
		t.Fatal("Количество публикация после раскодирования должно быть больше 0")
	}
}
