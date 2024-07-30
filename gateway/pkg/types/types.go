package types

// Публикация.
type Post struct {
	ID      int    `json:"id"`       // номер записи
	Title   string `json:"title"`    // заголовок публикации
	Content string `json:"content"`  // содержание публикации
	PubTime int64  `json:"pub_time"` // время публикации
	Link    string `json:"link"`     // ссылка на источник
}

type Comment struct {
	ID        int    `json:"id"`         // номер записи
	Post      int    `json:"post"`       // заголовок публикации
	Parent    int    `json:"parent"`     // заголовок публикации
	Message   string `json:"message"`    // содержание публикации
	CreatedAt int64  `json:"created_at"` // время публикации
}

type NewsShortDetailed struct {
	Post     Post      `json:"post"`
	Comments []Comment `json:"comments"`
}
