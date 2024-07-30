// Пакет для работы с RSS-потоками.
package rss

import (
	"aggregator/pkg/storage"
	"encoding/xml"
	"io"
	"net/http"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

// Лента записей RSS.
type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// Канал записей RSS.
type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

// Элемент записи RSS.
type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Link        string `xml:"link"`
}

// Parse принимат на вход url строку rss потока
// и возвращает раскодированные записи.
func Parse(url string) ([]storage.Post, error) {
	req, err := http.NewRequest("GET", url, nil)
	// обязательно проверка ошибки
	if err != nil {
		return nil, err
	}

	// для имитации реального запроса пользователя
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0")

	// запрос на получение XML rss потока
	client := &http.Client{}
	res, err := client.Do(req)
	// обязательно проверка ошибки
	if err != nil {
		return nil, err
	}

	// по завершению работы раскодировщика закрываем соединение
	defer res.Body.Close()

	// чтение сетевого ответа
	data, err := io.ReadAll(res.Body)
	// обязательно проверка ошибки
	if err != nil {
		return nil, err
	}

	// объявление переменной типа Feed для RSS потока
	var feed Feed

	// раскодирование полученного ответа в тип Feed
	err = xml.Unmarshal(data, &feed)
	// обязательно проверка ошибки
	if err != nil {
		return nil, err
	}

	// объявление слайса публикаций
	var news []storage.Post

	// утилита для очистки от тегов и XSS вставок
	p := bluemonday.StripTagsPolicy()

	// проходимся по записям в RSS ленте
	for _, item := range feed.Channel.Items {
		// попытка раскодирования даты публикации с учетом часового пояса
		date, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			// попытка раскодирования даты публикации без часового пояса
			date, err = time.Parse("Mon, 2 Jan 2006 15:04:05 GMT", item.PubDate)
			if err != nil {
				// если не получилось раскодировать дату публикции, то просто записываем 0 значение
				date = time.Unix(0, 0)
			}
		}

		// добавление в слайс публикаций
		news = append(news, storage.Post{
			Title:   item.Title,
			Content: p.Sanitize(item.Description),
			PubTime: date.Unix(),
			Link:    item.Link,
		})
	}

	return news, nil

}
