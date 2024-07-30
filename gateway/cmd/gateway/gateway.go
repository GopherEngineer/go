// Сервер "APIGateway".
package main

import (
	"fmt"
	"gateway/pkg/api"
	"log"
	"net/http"
)

// Порт сетевой службы
const PORT = "8080"

func main() {
	// создание API для обрабоки внешних HTTP запросов
	api := api.New()

	fmt.Println("HTTP server 'APIGateway' is started on localhost:" + PORT)
	// запуск сетевой службы
	log.Fatal(http.ListenAndServe(":"+PORT, api.Router()))
}
