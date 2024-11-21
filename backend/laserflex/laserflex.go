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

func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	// Открываем текущую директорию
	dirEntries, err := os.ReadDir(".")
	if err != nil {
		http.Error(w, "Failed to read directory", http.StatusInternalServerError)
		return
	}

	// Формируем HTML для отображения списка файлов
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<html><body><h1>Available Files for Download</h1><ul>"))

	for _, entry := range dirEntries {
		if strings.HasPrefix(entry.Name(), "file_downloaded_xls") && strings.HasSuffix(entry.Name(), ".xlsx") {
			fileName := entry.Name()
			downloadURL := fmt.Sprintf("/download/%s", fileName)
			w.Write([]byte(fmt.Sprintf("<li><a href=\"%s\">%s</a></li>", downloadURL, fileName)))
		}
	}

	w.Write([]byte("</ul></body></html>"))
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем имя файла из пути
	fileName := strings.TrimPrefix(r.URL.Path, "/download/")

	// Проверяем, существует ли файл
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Устанавливаем заголовки для скачивания
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// Открываем файл
	file, err := os.Open(fileName)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Отправляем файл в ответ
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Failed to send file", http.StatusInternalServerError)
		return
	}
}

func downloadFile(downloadURL string, iteration int) error {
	// Формируем имя файла
	fileName := fmt.Sprintf("file_downloaded_xls%d.xlsx", iteration)

	// Выполняем GET-запрос
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	// Создаём файл для сохранения
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Сохраняем данные в файл
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	log.Printf("File saved as: %s\n", fileName)
	return nil
}

var downloadCounter int

// Сюда нам прилетает файл из Битрикс24

func LaserflexGetFile(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection is starting...")

	contentType := r.Header.Get("Content-Type")
	log.Println("Content-Type:", contentType)

	var fileID string

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
		}
	}

	// Проверяем, есть ли необходимые параметры
	if fileID == "" {
		http.Error(w, "Missing file_id parameter", http.StatusBadRequest)
		return
	}

	// Вызов GetFileDetails для получения данных о файле
	fileDetails, err := GetFileDetails(fileID)
	if err != nil {
		log.Println("Error getting file details:", err)
		http.Error(w, "Failed to get file details", http.StatusInternalServerError)
		return
	}

	// Логируем DOWNLOAD_URL
	log.Printf("DOWNLOAD_URL: %s\n", fileDetails.DownloadURL)

	// Скачиваем файл с итерацией имени
	downloadCounter++
	err = downloadFile(fileDetails.DownloadURL, downloadCounter)
	if err != nil {
		log.Println("Error downloading file:", err)
		http.Error(w, "Failed to download file", http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("File downloaded successfully as file_downloaded_xls%d.xlsx", downloadCounter)))
}

func GetFileDetails(fileID string) (*FileDetails, error) {
	// Явно указываем URL с токеном
	clientEndpoint := "https://bitrix.laser-flex.ru/rest/149/ptosz34j8t6cpvgb/"
	requestURL := fmt.Sprintf("%sdisk.file.get.json?id=%s", clientEndpoint, fileID)

	// Создаём новый HTTP-запрос
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Println("Error creating new request:", err)
		return nil, err
	}

	// Устанавливаем Content-Type
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
