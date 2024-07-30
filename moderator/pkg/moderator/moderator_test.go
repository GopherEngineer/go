package moderator

import "testing"

func TestModerate(t *testing.T) {

	t.Run("without bad words", func(t *testing.T) {
		if !Moderate("The Golang is a super language!") {
			t.Fatal("Модерирование цензурной строки не пройдено")
		}
	})

	t.Run("with bad words", func(t *testing.T) {
		if Moderate("The best downloadable software for qwerty phones.") {
			t.Fatal("Модерирование нецензурной строки не пройдено")
		}
	})

}
