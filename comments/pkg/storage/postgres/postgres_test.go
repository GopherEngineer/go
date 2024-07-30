package postgres

import (
	"comments/pkg/storage"
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

		comment := storage.Comment{
			Post:      1,
			Parent:    0,
			Message:   "Текст комментария",
			CreatedAt: 0,
		}

		err = db.SaveComment(comment)
		if err != nil {
			t.Fatal(err)
		}

		comments, err := db.Comments(1)
		if err != nil {
			t.Fatal(err)
		}

		if len(comments) != 1 {
			t.Fatal("Сохраненных публикаций должно быть: 1")
		}
	})
}
