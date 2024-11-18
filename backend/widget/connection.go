package widget

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Ключ для хранения значения DealID в контексте
type key int

var DealIdGlobal string

var GlobalAuthIdWidget string

// Хендлер, который обрабатывает первый запрос и извлекает ID
func ConnectionBitrixWidget(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection is starting...")

	// Чтение тела запроса
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading request body:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Println("FULL REQUEST CONNECTION: ", string(bs))
	// Разбор строки запроса
	queryParams, err := url.ParseQuery(string(bs))
	if err != nil {
		log.Printf("Error parsing query: %v", err)
		http.Error(w, "Error parsing query", http.StatusBadRequest)
		return
	}

	// Получение и раскодирование параметра PLACEMENT_OPTIONS
	placementOptionsEncoded := queryParams.Get("PLACEMENT_OPTIONS")
	GlobalAuthIdWidget = queryParams.Get("AUTH_ID")
	log.Println("GlobalAuthIdWidget: ", GlobalAuthIdWidget)
	placementOptionsDecoded, err := url.QueryUnescape(placementOptionsEncoded)
	if err != nil {
		log.Printf("Error decoding PLACEMENT_OPTIONS: %v", err)
		http.Error(w, "Error decoding placement options", http.StatusBadRequest)
		return
	}
	log.Println("Decoded PLACEMENT_OPTIONS:", placementOptionsDecoded)

	// Парсинг JSON и извлечение поля ID
	var placementOptions map[string]interface{}
	if err := json.Unmarshal([]byte(placementOptionsDecoded), &placementOptions); err != nil {
		log.Printf("Error unmarshaling PLACEMENT_OPTIONS JSON: %v", err)
		http.Error(w, "Error unmarshaling placement options", http.StatusBadRequest)
		return
	}

	// Извлечение значения поля ID
	if id, ok := placementOptions["ID"].(string); ok {
		log.Printf("ID: %s", id)

		DealIdGlobal = id

		log.Println("Placement Options DealID Global", DealIdGlobal)
		// Проверяем, был ли ответ отправлен в SendDataForWidgetForm
		if w.Header().Get("Content-Type") == "" { // Если хедеры не установлены, редиректим
			redirectURL := "https://crmconsulting-api.ru/widget"
			log.Println("Redirecting to:", redirectURL)
			http.Redirect(w, r, redirectURL, http.StatusFound)
		}
	} else {
		log.Println("ID not found in PLACEMENT_OPTIONS")
		http.Error(w, "ID not found", http.StatusBadRequest)
	}
}
