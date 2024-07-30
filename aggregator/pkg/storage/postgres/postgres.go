// Пакет для работы с БД приложения "Новостной агрегатор".
package postgres

import (
	"aggregator/pkg/storage"
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

// News возвращает список публикаций из БД.
// Принимает два параметра: page - номер страницы, search - строка поиска в заголовке публикации.
func (s *Storage) News(page int, search string) ([]storage.Post, int, error) {
	// номер страницы не может быть меньше 1
	if page < 1 {
		page = 1
	}

	// переменная для сохранения количества публикаций в БД
	var count int
	// запрос на подчет количества публикаций с учетом поиска
	row := s.db.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM news
		WHERE title ILIKE $1;
	`, "%"+search+"%")
	err := row.Scan(&count)
	// обязательно проверяем на ошибку
	if err != nil {
		return nil, 0, err
	}

	// получаем публицации из базы данных с указанным ограничение
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			title,
			content,
			pub_time,
			link
		FROM news
		WHERE
			title ILIKE $2
		ORDER BY pub_time DESC
		LIMIT 15
		OFFSET $1;
	`, (page-1)*15, "%"+search+"%")
	// обязательно проверяем на ошибку
	if err != nil {
		return nil, 0, err
	}

	// объявляем слайс для публикаций
	news := make([]storage.Post, 0)

	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var p storage.Post

		// сохраняем полученные значение публикации в переменную
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, 0, err
		}

		// добавляем переменной в массив результатов
		news = append(news, p)
	}

	// ВАЖНО не забыть проверить rows.Err()
	return news, count, rows.Err()
}

// NewsDetail возвращает публикацию по ID из БД.
func (s *Storage) NewsDetail(id int) (*storage.Post, error) {
	// получаем публицацию из базы данных
	row := s.db.QueryRow(context.Background(), `
		SELECT 
			id,
			title,
			content,
			pub_time,
			link
		FROM news
		WHERE id = $1;
	`, id)

	var p storage.Post

	// сканирование полученной строки в переменную
	err := row.Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.PubTime,
		&p.Link,
	)
	// обязательно проверяем на ошибку
	if err != nil {
		return nil, err
	}

	// возвращаем полученную публикацию
	return &p, nil
}

// SaveNews сохраняет массив публикаций в базу данных.
// Уникальность публикации проверяется по полю "Link".
func (s *Storage) SaveNews(news []storage.Post) error {
	// проходим по слайсу публикаций и создаем новую публикацию в базе данных
	for _, post := range news {
		_, err := s.db.Exec(context.Background(), `
			INSERT INTO news (title, content, pub_time, link)
			VALUES ($1, $2, $3, $4);
		`,
			post.Title,
			post.Content,
			post.PubTime,
			post.Link,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
