// Сервер "Новостной агрегатор".
package main

import (
	"aggregator/pkg/api"
	"aggregator/pkg/rss"
	"aggregator/pkg/storage"
	"aggregator/pkg/storage/postgres"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// конфигурация приложения
type config struct {
	URLS   []string `json:"rss"`
	Period int      `json:"period"`
}

// порт сетевой службы
const PORT = "8081"

func main() {

	// подключение к базе данных где будут храниться публикации
	db, err := postgres.New()
	// обязательно обработка ошибки
	if err != nil {
		log.Fatal(err)
	}

	// создание API для обрабоки внешних HTTP запросов
	api := api.New(db)

	// чтение конфигурационного файла для получение
	// списка rss потоков и тайминга их обходов.
	// файл располагается в корне проекта
	b, err := os.ReadFile("./config.json")
	// обязательно обработка ошибки
	if err != nil {
		log.Fatal(err)
	}

	// раскодирование прочитанного файла конфигрурации
	var config config
	err = json.Unmarshal(b, &config)
	// обязательно обработка ошибки
	if err != nil {
		log.Fatal(err)
	}

	ch_news := make(chan []storage.Post)
	ch_errs := make(chan error)

	// запуск горутины чтения из канала публикаций
	// и сохранения их в базу данных
	go func() {
		for news := range ch_news {
			db.SaveNews(news)
		}
	}()

	// запуск горутины чтения из канала ошибок
	// и печати из в логи
	go func() {
		for err := range ch_errs {
			log.Println("Ошибка чтения RSS ленты:", err)
		}
	}()

	// чтение rss ссылок из массива
	// и запуск их обработки в горутине
	for _, url := range config.URLS {
		go func(url string) {
			for {

				// вызов rss парсера и получение
				// декодированных публикаций
				news, err := rss.Parse(url)
				if err != nil {
					// отправка ошибки в канал для далнейшей её обработки
					ch_errs <- err
					goto WAIT_PERIOD
				}

				// отправка полученных публикаций в канал
				ch_news <- news

			WAIT_PERIOD:
				// задержка на указанные период в конфигурации
				// чтобы не получить блокировку от rss источника
				time.Sleep(time.Minute * time.Duration(config.Period))

			}
		}(url)
	}

	fmt.Println("HTTP server 'Aggregator' is started on localhost:" + PORT)
	// запуск сетевой службы
	log.Fatal(http.ListenAndServe(":"+PORT, api.Router()))

}
