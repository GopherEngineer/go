// Сервер "Комментарии".
package main

import (
	"comments/pkg/api"
	"comments/pkg/storage/postgres"
	"fmt"
	"log"
	"net/http"
)

// Порт сетевой службы
const PORT = "8082"

func main() {

	// подключение к базе данных где будут храниться комментарии
	db, err := postgres.New()
	// обязательно обработка ошибки
	if err != nil {
		log.Fatal(err)
	}

	// создание API для обрабоки внешних HTTP запросов
	api := api.New(db)

	fmt.Println("HTTP server 'Comments' is started on localhost:" + PORT)
	// запуск сетевой службы
	log.Fatal(http.ListenAndServe(":"+PORT, api.Router()))
}
