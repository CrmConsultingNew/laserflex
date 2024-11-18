package laserflex

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type FileDetails struct {
	ID          string `json:"ID"`
	Name        string `json:"NAME"`
	DownloadURL string `json:"DOWNLOAD_URL"`
}

func LaserflexGetFile(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection is starting...")

	contentType := r.Header.Get("Content-Type")
	log.Println("Content-Type:", contentType)

	var fileID, authToken, clientEndpoint string

	// Извлекаем параметры из запроса
	switch {
	case strings.HasPrefix(contentType, "multipart/form-data"):
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			log.Println("Error parsing multipart data:", err)
			http.Error(w, "Invalid multipart data", http.StatusBadRequest)
			return
		}
		for key, values := range r.Form {
			log.Printf("Form Field %s: %v\n", key, values)
			if key == "file_id" && len(values) > 0 {
				fileID = values[0]
			}
			if key == "auth[client_endpoint]" && len(values) > 0 {
				clientEndpoint = values[0]
			}
			if key == "auth[member_id]" && len(values) > 0 {
				authToken = values[0]
			}
		}

	case contentType == "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			log.Println("Error parsing URL-encoded data:", err)
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		for key, values := range r.Form {
			log.Printf("Form Field %s: %v\n", key, values)
			if key == "file_id" && len(values) > 0 {
				fileID = values[0]
			}
			if key == "auth[client_endpoint]" && len(values) > 0 {
				clientEndpoint = values[0]
			}
			if key == "auth[member_id]" && len(values) > 0 {
				authToken = values[0]
			}
		}
	}

	// Проверяем, есть ли необходимые параметры
	if fileID == "" || clientEndpoint == "" || authToken == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Вызов метода disk.file.get
	fileDetails, err := GetFileDetails(fileID, authToken, clientEndpoint)
	if err != nil {
		log.Println("Error getting file details:", err)
		http.Error(w, "Failed to get file details", http.StatusInternalServerError)
		return
	}

	// Логируем DOWNLOAD_URL
	log.Printf("DOWNLOAD_URL: %s\n", fileDetails.DownloadURL)

	// Отправляем ответ клиенту
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("File Download URL: %s", fileDetails.DownloadURL)))
}

func GetFileDetails(fileID string, authToken string, clientEndpoint string) (*FileDetails, error) {
	// Формируем URL для запроса
	bitrixMethod := "disk.file.get.json"
	requestURL := fmt.Sprintf("%s%s?id=%s", clientEndpoint, bitrixMethod, fileID)

	// Создаём новый HTTP-запрос
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Println("Error creating new request:", err)
		return nil, err
	}

	// Добавляем заголовки авторизации
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: received status code %d\n", resp.StatusCode)
		return nil, fmt.Errorf("failed to get file details: status %d", resp.StatusCode)
	}

	// Читаем тело ответа
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil, err
	}

	// Логируем ответ для отладки
	log.Println("GetFileDetails Response:", string(responseData))

	// Парсим ответ в структуру FileDetails
	var response struct {
		Result FileDetails `json:"result"`
	}
	if err := json.Unmarshal(responseData, &response); err != nil {
		log.Println("Error unmarshaling response:", err)
		return nil, err
	}

	return &response.Result, nil
}
