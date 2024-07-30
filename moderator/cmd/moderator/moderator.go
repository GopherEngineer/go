// Сервер "Модератор комментариев".
package main

import (
	"fmt"
	"log"
	"moderator/pkg/api"
	"net/http"
)

// Порт сетевой службы
const PORT = "8083"

func main() {
	// создание API для обрабоки внешних HTTP запросов
	api := api.New()

	fmt.Println("HTTP server 'Moderator' is started on localhost:" + PORT)
	// запуск сетевой службы
	log.Fatal(http.ListenAndServe(":"+PORT, api.Router()))
}
