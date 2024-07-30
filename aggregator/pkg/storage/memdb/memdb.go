// Заглушка базы данных для тестирования.
package memdb

import (
	"aggregator/pkg/storage"
	"strconv"
)

// Хранилище данных.
type Storage struct{}

// Конструктор, принимает строку подключения к БД.
func New() (*Storage, error) {
	return new(Storage), nil
}

// News возвращает список публикаций.
func (s *Storage) News(page int, search string) ([]storage.Post, int, error) {
	var news []storage.Post

	for i := range 15 {
		news = append(news, storage.Post{
			ID:      i,
			Title:   "Title",
			Content: "Content",
			Link:    strconv.Itoa(i),
		})
	}

	return news, 15, nil
}

// NewsDetail возвращает публикацию по ID.
func (s *Storage) NewsDetail(id int) (*storage.Post, error) {
	return &storage.Post{
		ID:      id,
		Title:   "Title",
		Content: "Content",
		Link:    strconv.Itoa(id),
	}, nil
}

// SaveNews сохраняет массив публикаций в базу данных.
// Уникальность публикации проверяется по полю "Link".
func (s *Storage) SaveNews(news []storage.Post) error {
	return nil
}
