package postgres

import (
	"aggregator/pkg/storage"
	"context"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestDB(t *testing.T) {
	ctx := context.Background()

	name := "news"
	user := "user"
	password := "password"

	container, err := postgres.Run(
		ctx,
		"docker.io/postgres:16-alpine",
		postgres.WithDatabase(name),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	if err != nil {
		log.Fatalf("Не удалось запустить контейнер: %s", err)
	}

	b, err := os.ReadFile("../../../schema.sql")
	if err != nil {
		log.Fatalf("Не удалось прочитать файл schema.sql: %s", err)
	}

	// миграция схемы таблицы
	_, _, err = container.Exec(ctx, []string{"psql", "-U", user, "-d", name, "-c", string(b)})
	if err != nil {
		t.Fatal(err)
	}

	// создание снимка базы данных
	// для переиспользования начального состояния для тестов
	err = container.Snapshot(ctx, postgres.WithSnapshotName("snapshot"))
	if err != nil {
		t.Fatal(err)
	}

	// очиститка конейнера после завершения тестирования
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Не удалось очистить контейнер: %s", err)
		}
	})

	// получение ссылки на подключение к базе данных
	dbURL, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("DB_URL", dbURL)

	t.Run("empty_db_url", func(t *testing.T) {
		t.Setenv("DB_URL", "")
		if _, err := New(); err == nil {
			t.Fatal("Нет проверки на пустую переменную окружения DB_URL")
		}
	})

	t.Run("new", func(t *testing.T) {
		storage, err := New()
		if err != nil {
			t.Fatal(err)
		}
		storage.db.Close()
	})

	t.Run("save_and_fetch", func(t *testing.T) {
		t.Cleanup(func() {
			if err := container.Restore(ctx); err != nil {
				t.Fatal(err)
			}
		})

		db, err := New()
		if err != nil {
			t.Fatal(err)
		}

		posts := []storage.Post{
			{
				Title:   "Title 1",
				Content: "Content 1",
				PubTime: 0,
				Link:    "Link 1",
			},
			{
				Title:   "Title 2",
				Content: "Content 2",
				PubTime: 0,
				Link:    "Link 2",
			},
		}

		err = db.SaveNews(posts)
		if err != nil {
			t.Fatal(err)
		}

		// передача 0 для провеки, что будет значение по умолчанию 10
		news, _, err := db.News(0, "")
		if err != nil {
			t.Fatal(err)
		}

		if len(news) != 2 {
			t.Fatal("Сохраненных публикаций должно быть: 2")
		}
	})

	t.Run("uniq_links", func(t *testing.T) {
		t.Cleanup(func() {
			if err := container.Restore(ctx); err != nil {
				t.Fatal(err)
			}
		})

		db, err := New()
		if err != nil {
			t.Fatal(err)
		}

		posts := []storage.Post{
			{
				Title:   "Title 1",
				Content: "Content 1",
				PubTime: 0,
				Link:    "Link 1",
			},
			{
				Title:   "Title 2",
				Content: "Content 2",
				PubTime: 0,
				Link:    "Link 1",
			},
		}

		err = db.SaveNews(posts)
		if err == nil {
			t.Fatal("Нет проверки на уникальность ссылок публикаций")
		}
	})
}
