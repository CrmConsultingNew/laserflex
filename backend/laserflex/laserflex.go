package laserflex

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func LaserflexGetFile(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection is starting...")

	// Чтение тела запроса
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Логируем полный запрос
	log.Println("FULL REQUEST CONNECTION: ", string(bs))

	// Парсинг данных
	data, err := url.ParseQuery(string(bs))
	if err != nil {
		log.Println("Error parsing request body:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Пример извлечения данных
	log.Println("Parsed Data:")
	for key, values := range data {
		log.Printf("%s: %s\n", key, values)
	}

	// Пример использования значений
	documentID := data["document_id[]"]
	authDomain := data.Get("auth[domain]") // Извлечение одного значения

	log.Println("Document ID:", documentID)
	log.Println("Auth Domain:", authDomain)

	// Ответ клиенту
	fmt.Fprintln(w, "Data parsed successfully")
}
