package laserflex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func LaserflexGetFile(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection is starting...")

	contentType := r.Header.Get("Content-Type")
	log.Println("Content-Type:", contentType)

	var fileID, dealID, smartProcessID string

	// Извлекаем параметры из URL
	queryParams := r.URL.Query()

	// Считываем необходимые параметры
	dealID = queryParams.Get("deal_id")
	smartProcessID = queryParams.Get("smartProcessID")
	fileID = queryParams.Get("file_id")

	if dealID != "" {
		log.Printf("Extracted deal_id: %s\n", dealID)
	}

	if smartProcessID != "" {
		log.Printf("Extracted smartProcessID: %s\n", smartProcessID)
	}

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
	fileName := fmt.Sprintf("file_downloaded_xls%d.xlsx", downloadCounter)
	err = downloadFile(fileDetails.DownloadURL, downloadCounter)
	if err != nil {
		log.Println("Error downloading file:", err)
		http.Error(w, "Failed to download file", http.StatusInternalServerError)
		return
	}

	// Чтение продуктов из Excel файла
	products, err := ReadXlsProductRows(fileName)
	if err != nil {
		log.Println("Error reading Excel file:", err)
		http.Error(w, "Failed to process Excel file", http.StatusInternalServerError)
		return
	}

	// Создаем массив для хранения ID товаров
	var productIDs []int

	// Добавление продуктов в Bitrix24
	for _, product := range products {
		productID, err := AddProductsWithImage(product, "52") // Используем ID раздела "52" как пример
		if err != nil {
			log.Printf("Error adding product %s: %v", product.Name, err)
			continue
		}
		productIDs = append(productIDs, productID)
	}

	log.Printf("Processed products: %+v\n", products)
	log.Printf("Added Product IDs: %+v\n", productIDs)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File processed and products added successfully"))
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

func AddProductsRowToDeal(name string, currency string, price float64, sort int, authID string) error {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "crm.deal.productrows.set"

	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	requestBody := map[string]interface{}{
		"fields": map[string]interface{}{
			"NAME":        name,
			"CURRENCY_ID": currency,
			"PRICE":       price,
			"SORT":        sort,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(responseData, &response); err != nil {
		return err
	}

	if _, ok := response["error"]; ok {
		return fmt.Errorf("Ошибка: %s", response["error_description"])
	}

	log.Println("Товар добавлен с ID:", response["result"])
	return nil
}
