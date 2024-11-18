package laserflex

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func LaserflexGetFile(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection is starting...")

	contentType := r.Header.Get("Content-Type")
	log.Println("Content-Type:", contentType)

	var fileURL string

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
			if key == "doc" && len(values) > 0 {
				fileURL = values[0]
			}
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
			if key == "doc" && len(values) > 0 {
				fileURL = values[0]
			}
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

	// Проверяем, есть ли ссылка на файл
	if fileURL == "" {
		http.Error(w, "Missing file URL", http.StatusBadRequest)
		return
	}

	log.Println("File URL:", fileURL)

	// Скачиваем файл
	err := downloadFile(fileURL, "downloaded_file.xlsx")
	if err != nil {
		log.Println("Error downloading file:", err)
		http.Error(w, "Failed to download file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File downloaded successfully."))
}

// Функция для скачивания файла
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
