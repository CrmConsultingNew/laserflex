package laserflex

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type FileDetails struct {
	ID          string `json:"ID"`
	Name        string `json:"NAME"`
	DownloadURL string `json:"DOWNLOAD_URL"`
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

	// Проверяем код ответа
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

	// Логируем необработанный ответ для отладки
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
			if key == "auth[member_id]" && len(values) > 0 { // Пример для получения токена
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
			if key == "auth[member_id]" && len(values) > 0 { // Пример для получения токена
				authToken = values[0]
			}
		}
	}

	// Проверяем, есть ли file_id и другие параметры
	if fileID == "" || clientEndpoint == "" || authToken == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Получаем информацию о файле
	fileDetails, err := GetFileDetails(fileID, authToken, clientEndpoint)
	if err != nil {
		log.Println("Error getting file details:", err)
		http.Error(w, "Failed to get file details", http.StatusInternalServerError)
		return
	}

	// Логируем детали файла
	log.Printf("File Details: %+v\n", fileDetails)

	// Пример: Скачать файл с помощью URL
	err = downloadFile(fileDetails.DownloadURL, fileDetails.Name)
	if err != nil {
		log.Println("Error downloading file:", err)
		http.Error(w, "Failed to download file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("File '%s' downloaded successfully.", fileDetails.Name)))
}

func downloadFile(fileURL, filePath string) error {
	// Выполняем GET-запрос
	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("failed to fetch file: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch file: status %d", resp.StatusCode)
	}

	// Создаём файл для сохранения
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Сохраняем данные в файл
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	log.Printf("File downloaded to %s", filePath)
	return nil
}
