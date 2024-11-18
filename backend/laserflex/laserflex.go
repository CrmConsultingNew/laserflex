package laserflex

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
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

	// Определение Content-Type
	contentType := r.Header.Get("Content-Type")
	log.Println("Content-Type:", contentType)

	// Обработка JSON
	if contentType == "application/json" {
		var data map[string]interface{}
		if err := json.Unmarshal(bs, &data); err != nil {
			log.Println("Error parsing JSON:", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		log.Println("Parsed JSON:", data)
	}

	// Обработка URL-encoded данных
	if contentType == "application/x-www-form-urlencoded" {
		parsedData, err := url.ParseQuery(string(bs))
		if err != nil {
			log.Println("Error parsing URL-encoded data:", err)
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}
		log.Println("Parsed Form Data:")
		for key, values := range parsedData {
			log.Printf("%s: %s\n", key, values)
		}
	}

	// Обработка multipart данных
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
			log.Println("Error parsing multipart data:", err)
			http.Error(w, "Invalid multipart data", http.StatusBadRequest)
			return
		}
		log.Println("Parsed Multipart Form:")
		for key, values := range r.Form {
			log.Printf("%s: %s\n", key, values)
		}
	}

	// Если Content-Type неизвестен
	if contentType == "" {
		log.Println("Raw Body Data:", string(bs))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request processed successfully."))
}
