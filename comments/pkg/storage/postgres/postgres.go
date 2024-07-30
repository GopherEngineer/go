// Пакет для работы с БД сервера "Комментарии".
package postgres

import (
	"comments/pkg/storage"
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func New() (*Storage, error) {
	url := os.Getenv("DB_URL")
	if url == "" {
		return nil, errors.New("переменная окружения DB_URL не задана")
	}
	db, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// Comments возвращает список комментариев для публикаций из БД.
// Принимает один параметр id - идентификатор публикации.
func (s *Storage) Comments(id int) ([]storage.Comment, error) {

	// получаем комментарии из базы данных
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			post,
			parent,
			message,
			created_at
		FROM comments
		WHERE post = $1
		ORDER BY created_at ASC;
	`, id)
	// обязательно проверяем на ошибку
	if err != nil {
		return nil, err
	}

	// объявляем слайс для комментариев
	comments := []storage.Comment{}

	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var c storage.Comment

		// сохраняем полученные значение комментария в переменную
		err = rows.Scan(
			&c.ID,
			&c.Post,
			&c.Parent,
			&c.Message,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// добавляем переменной в массив результатов
		comments = append(comments, c)
	}

	// ВАЖНО не забыть проверить rows.Err()
	return comments, rows.Err()
}

// SaveComment сохраняет комментарий для публикации в базу данных.
func (s *Storage) SaveComment(c storage.Comment) error {
	_, err := s.db.Exec(context.Background(), `
			INSERT INTO comments (post, parent, message, created_at)
			VALUES ($1, $2, $3, $4);
		`,
		c.Post,
		c.Parent,
		c.Message,
		c.CreatedAt,
	)
	return err
}
