package laserflex

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func LaserflexGetFile(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection is starting...")

	contentType := r.Header.Get("Content-Type")
	log.Println("Content-Type:", contentType)

	switch {
	case strings.HasPrefix(contentType, "multipart/form-data"):
		// Обработка multipart данных
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			log.Println("Error parsing multipart data:", err)
			http.Error(w, "Invalid multipart data", http.StatusBadRequest)
			return
		}
		for key, values := range r.Form {
			log.Printf("Form Field %s: %v\n", key, values)
		}

	case contentType == "application/x-www-form-urlencoded":
		// Обработка URL-encoded данных
		if err := r.ParseForm(); err != nil {
			log.Println("Error parsing URL-encoded data:", err)
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		for key, values := range r.Form {
			log.Printf("Form Field %s: %v\n", key, values)
		}

	default:
		// Чтение тела запроса для других форматов
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading request body:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		log.Println("Raw Body:", string(body))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request processed successfully."))
}
