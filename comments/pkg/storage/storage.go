// Пакет для соблюдения контракта работы с базой данных.
package storage

// Комментарий.
type Comment struct {
	ID        int    `json:"id"`         // номер записи
	Post      int    `json:"post"`       // ID публикации
	Parent    int    `json:"parent"`     // ID родителького комментария
	Message   string `json:"message"`    // текст комментария
	CreatedAt int64  `json:"created_at"` // время публикации
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Comments(n int) ([]Comment, error)
	SaveComment(Comment) error
}
