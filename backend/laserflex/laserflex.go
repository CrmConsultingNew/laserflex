package laserflex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func LaserflexGetFile(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection is starting...")

	contentType := r.Header.Get("Content-Type")
	log.Println("Content-Type:", contentType)

	var fileID, dealID, smartProcessID string
	var assignedById int
	//var docIDs []int // Массив для хранения docId

	// Извлекаем параметры из URL
	queryParams := r.URL.Query()

	// Считываем необходимые параметры
	dealID = queryParams.Get("deal_id")
	smartProcessID = queryParams.Get("smartProcessID")
	fileID = queryParams.Get("file_id")
	assignedByIdStr := queryParams.Get("assignedById")

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

	// Преобразование assignedById в int
	if assignedByIdStr != "" {
		var err error
		assignedById, err = strconv.Atoi(assignedByIdStr)
		if err != nil {
			log.Printf("Error converting assignedById to int: %v\n", err)
			http.Error(w, "Invalid assignedById parameter", http.StatusBadRequest)
			return
		}
		log.Printf("Extracted assignedById: %d\n", assignedById)
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

	/*// Чтение продуктов из Excel файла
	products, err := ReadXlsProductRows(fileName)
	if err != nil {
		log.Println("Error reading Excel file:", err)
		http.Error(w, "Failed to process Excel file", http.StatusInternalServerError)
		return
	}*/

	// Создаем массив для хранения ID товаров
	/*var productIDs []int
	var totalProductsPrice float64

	// Добавление продуктов в Bitrix24
	for _, product := range products {
		productID, err := AddProductsWithImage(product, "52") // Используем ID раздела "52" как пример
		if err != nil {
			log.Printf("Error adding product %s: %v", product.Name, err)
			continue
		}
		productIDs = append(productIDs, productID)
		totalProductsPrice += product.Price * product.Quantity // Учитываем общую цену с учетом количества
	}*/

	// После получения productIDs и products
	/*var quantities []float64
	var prices []float64
	for _, product := range products {
		quantities = append(quantities, product.Quantity)
		prices = append(prices, product.Price)
	}

	// Добавление товаров в сделку
	err = AddProductsRowToDeal(dealID, productIDs, quantities, prices)
	if err != nil {
		log.Printf("Error adding product rows to deal: %v", err)
		http.Error(w, "Failed to add product rows to deal", http.StatusInternalServerError)
		return
	}

	// Добавление документа в Bitrix24
	docId, err := AddCatalogDocument(dealID, assignedById, totalProductsPrice)
	if err != nil {
		log.Printf("Error adding catalog document: %v", err)
		http.Error(w, "Failed to add catalog document", http.StatusInternalServerError)
		return
	}

	// Добавляем docId в массив
	docIDs = append(docIDs, docId)

	if len(productIDs) != len(quantities) {
		log.Println("Mismatched lengths: productIDs and quantities")
		http.Error(w, "Mismatched lengths of productIDs and quantities", http.StatusInternalServerError)
		return
	}

	for i, productId := range productIDs {
		quantity := quantities[i]

		err := AddCatalogDocumentElement(docId, productId, quantity) // добавить товары в документ прихода
		if err != nil {
			log.Printf("Error adding catalog document with element: %v", err)
			http.Error(w, "Failed to add catalog document with element", http.StatusInternalServerError)
			return
		}
	}

	// Проведение документа
	err = ConductDocumentId(docId)
	if err != nil {
		log.Printf("Error conducting document: %v", err)
		http.Error(w, "Failed to conduct document", http.StatusInternalServerError)
		return
	}

	// Сохраняем docIDs в текстовый файл
	err = saveDocIDsToFile("document_ids.txt", docIDs)
	if err != nil {
		log.Printf("Error saving document IDs to file: %v", err)
		http.Error(w, "Failed to save document IDs to file", http.StatusInternalServerError)
		return
	}*/

	/*productionEngineerId, err := GetProductionEngineerIdByDeal(dealID)
	if err != nil {
		log.Printf("Error getting production engineer ID: %v", err)
	}*/

	// Конвертируем файл Excel в JSON
	StartJsonConverterFromExcel(fileName)

	// Читаем сгенерированный JSON
	var parsedData []ParsedData
	jsonData, err := os.ReadFile("output.json")
	if err != nil {
		log.Printf("Error reading JSON file: %v\n", err)
		http.Error(w, "Failed to read parsed JSON data", http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(jsonData, &parsedData); err != nil {
		log.Printf("Error unmarshalling JSON data: %v\n", err)
		http.Error(w, "Failed to parse JSON data", http.StatusInternalServerError)
		return
	}

	// Лимит итераций
	iterationLimit := 10
	iterationCount := 0

	// Обрабатываем каждый блок данных из JSON
	for _, data := range parsedData {
		if iterationCount >= iterationLimit {
			log.Println("Reached iteration limit of 10, stopping further processing.")
			break
		}

		if data.LaserWorks != nil {
			taskID, err := AddTaskToGroup("laser_works", 149, data.LaserWorks.GroupID, 1046, 458)
			if err != nil {
				log.Printf("Error creating laser_works task: %v\n", err)
				continue
			}
			log.Printf("LaserWorks Task created with ID: %d\n", taskID)

			// Создаем подзадачи для laser_works
			customFields := CustomTaskFields{
				OrderNumber: data.LaserWorks.Data["№ заказа"],
				Customer:    data.LaserWorks.Data["Заказчик"],
				Manager:     data.LaserWorks.Data["Менеджер"],
				Material:    data.LaserWorks.Data["Количество материала"],
				Comment:     data.LaserWorks.Data["Комментарий"],
				Coating:     data.LaserWorks.Data["Нанесение покрытий"],
			}
			subTaskID, err := AddTaskToParentId("laser_works_subtask", 149, data.LaserWorks.GroupID, taskID, customFields)
			if err != nil {
				log.Printf("Error creating laser_works subtask: %v\n", err)
				continue
			}
			log.Printf("LaserWorks Subtask created with ID: %d\n", subTaskID)
		}

		if data.BendWorks != nil {
			taskID, err := AddTaskToGroup("bend_works", 149, data.BendWorks.GroupID, 1046, 458)
			if err != nil {
				log.Printf("Error creating bend_works task: %v\n", err)
				continue
			}
			log.Printf("BendWorks Task created with ID: %d\n", taskID)

			// Создаем подзадачи для bend_works
			customFields := CustomTaskFields{
				OrderNumber: data.BendWorks.Data["№ заказа"],
				Customer:    data.BendWorks.Data["Заказчик"],
				Manager:     data.BendWorks.Data["Менеджер"],
				Material:    data.BendWorks.Data["Количество материала"],
				Comment:     data.BendWorks.Data["Комментарий"],
				Coating:     data.BendWorks.Data["Нанесение покрытий"],
			}
			subTaskID, err := AddTaskToParentId("bend_works_subtask", 149, data.BendWorks.GroupID, taskID, customFields)
			if err != nil {
				log.Printf("Error creating bend_works subtask: %v\n", err)
				continue
			}
			log.Printf("BendWorks Subtask created with ID: %d\n", subTaskID)
		}

		if data.Production != nil {
			checklist := []map[string]interface{}{}
			for _, step := range strings.Fields(data.Production.Data["Производство"]) {
				checklist = append(checklist, map[string]interface{}{
					"title": step,
				})
			}

			taskID, err := AddTaskWithChecklist("production", 149, 1046, 458, checklist)
			if err != nil {
				log.Printf("Error creating production task: %v\n", err)
				continue
			}
			log.Printf("Production Task with checklist created with ID: %d\n", taskID)
		}

		iterationCount++
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File processed and products added successfully"))
}

// Функция для сохранения docIDs в текстовый файл
func saveDocIDsToFile(filename string, docIDs []int) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	for _, id := range docIDs {
		_, err := file.WriteString(fmt.Sprintf("%d\n", id))
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	log.Printf("Document IDs saved to file: %s", filename)
	return nil
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

func AddProductsRowToDeal(dealID string, productIDs []int, quantities []float64, prices []float64) error {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "crm.deal.productrows.set"

	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	// Создаем массив строк для отправки
	var rows []map[string]interface{}
	for i, productID := range productIDs {
		rows = append(rows, map[string]interface{}{
			"PRODUCT_ID": productID,
			"QUANTITY":   quantities[i],
			"PRICE":      prices[i],
		})
	}

	requestBody := map[string]interface{}{
		"id":   dealID,
		"rows": rows,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %v", err)
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(responseData, &response); err != nil {
		return fmt.Errorf("error unmarshalling response: %v", err)
	}

	if _, ok := response["error"]; ok {
		return fmt.Errorf("Ошибка: %s", response["error_description"])
	}

	log.Println("Товарные строки добавлены в сделку:", dealID)
	return nil
}
