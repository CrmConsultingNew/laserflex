package laserflex

import (
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
