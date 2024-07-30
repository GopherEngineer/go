// Пакет для соблюдения контракта работы с базой данных.
package storage

// Публикация, получаемая из RSS.
type Post struct {
	ID      int    `json:"id"`       // номер записи
	Title   string `json:"title"`    // заголовок публикации
	Content string `json:"content"`  // содержание публикации
	PubTime int64  `json:"pub_time"` // время публикации
	Link    string `json:"link"`     // ссылка на источник
}

// Публикации объедененные с пагинацией
type NewsFullDetailed struct {
	News  []Post `json:"news"`
	Page  int    `json:"page"`
	Pages int    `json:"pages"`
	Count int    `json:"count"`
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	News(n int, s string) ([]Post, int, error)
	NewsDetail(id int) (*Post, error)
	SaveNews([]Post) error
}
