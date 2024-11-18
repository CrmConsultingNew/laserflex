package laserflex

import (
	"log"
	"net/http"
)

func LaserflexGetFile(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection is starting...")

	// Извлечение параметров из строки запроса
	query := r.URL.Query()

	// Читаем параметры
	doc := query.Get("doc")                       // Параметр doc (ссылка на файл)
	userId := query.Get("userId")                 // Параметр userId
	clientID := query.Get("client_id")            // Параметр client_id
	clientSecret := query.Get("client_secret")    // Параметр client_secret
	smartProcessID := query.Get("smartProcessID") // Параметр smartProcessID

	// Логирование параметров
	log.Println("Parsed Query Params:")
	log.Printf("doc: %s\n", doc)
	log.Printf("userId: %s\n", userId)
	log.Printf("client_id: %s\n", clientID)
	log.Printf("client_secret: %s\n", clientSecret)
	log.Printf("smartProcessID: %s\n", smartProcessID)

	// Ответ сервером
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Query parameters received and processed successfully."))
}
