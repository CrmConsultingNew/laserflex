package laserflex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func LaserflexGetFile(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection is starting...")

	// Извлекаем параметры из URL
	queryParams := r.URL.Query()
	fileID := queryParams.Get("file_id")
	smartProcessIDStr := queryParams.Get("smartProcessID")
	engineerIDStr := queryParams.Get("engineer_id")
	dealID := queryParams.Get("deal_id")
	assignedByIdStr := queryParams.Get("assigned")
	assignedById, err := strconv.Atoi(assignedByIdStr)
	if err != nil {
		log.Printf("Error converting engineerID to int: %v\n", err)
		http.Error(w, "Invalid engineerID parameter", http.StatusBadRequest)
		return
	}
	engineerID, err := strconv.Atoi(engineerIDStr)
	if err != nil {
		log.Printf("Error converting engineerID to int: %v\n", err)
		http.Error(w, "Invalid engineerID parameter", http.StatusBadRequest)
		return
	}
	log.Printf("Engineer ID: %v", engineerID)

	if fileID == "" {
		http.Error(w, "Missing file_id parameter", http.StatusBadRequest)
		return
	}

	// Конвертация smartProcessID в int
	smartProcessID, err := strconv.Atoi(smartProcessIDStr)
	if err != nil {
		log.Printf("Error converting smartProcessID to int: %v\n", err)
		http.Error(w, "Invalid smartProcessID parameter", http.StatusBadRequest)
		return
	}

	// Получаем данные о файле
	fileDetails, err := GetFileDetails(fileID)
	if err != nil {
		log.Printf("Error getting file details: %v\n", err)
		http.Error(w, "Failed to get file details", http.StatusInternalServerError)
		return
	}

	// Скачиваем файл
	fileName := fmt.Sprintf("file_downloaded_xls%d.xlsx", downloadCounter)
	err = downloadFile(fileDetails.DownloadURL, downloadCounter)
	if err != nil {
		log.Printf("Error downloading file: %v\n", err)
		http.Error(w, "Failed to download file", http.StatusInternalServerError)
		return
	}

	products, err := ReadXlsProductRows(fileName)
	if err != nil {
		log.Println("Error reading Excel file:", err)
		http.Error(w, "Failed to process Excel file", http.StatusInternalServerError)
		return
	}

	var productIDs []int
	var totalProductsPrice float64

	for _, product := range products {
		productID, err := AddProductsWithImage(product, "52") // Используем ID раздела "52" как пример
		if err != nil {
			log.Printf("Error adding product %s: %v", product.Name, err)
			continue
		}
		productIDs = append(productIDs, productID)
		totalProductsPrice += product.Price * product.Quantity // Учитываем общую цену с учетом количества
	}

	var quantities []float64
	var prices []float64
	for _, product := range products {
		quantities = append(quantities, product.Quantity)
		prices = append(prices, product.Price)
	}

	err = AddProductsRowToDeal(dealID, productIDs, quantities, prices)
	if err != nil {
		log.Printf("Error adding product rows to deal: %v", err)
		http.Error(w, "Failed to add product rows to deal", http.StatusInternalServerError)
		return
	}

	docId, err := AddCatalogDocument(dealID, assignedById, totalProductsPrice)
	if err != nil {
		log.Printf("Error adding catalog document: %v", err)
		http.Error(w, "Failed to add catalog document", http.StatusInternalServerError)
		return
	}

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

	err = ConductDocumentId(docId)
	if err != nil {
		log.Printf("Error conducting document: %v", err)
		http.Error(w, "Failed to conduct document", http.StatusInternalServerError)
		return
	}

	var arrayOfTasksIDs []int

	// Обрабатываем задачи
	taskIDLaserWorks, err := processLaserWorks(fileName, smartProcessID, engineerID)
	if err != nil {
		log.Printf("Error processing Laser Works: %v\n", err)
		http.Error(w, "Failed to process Laser Works", http.StatusInternalServerError)
		return
	}

	taskIDBendWorks, err := processBendWorks(fileName, smartProcessID, engineerID)
	if err != nil {
		log.Printf("Error processing Bend Works: %v\n", err)
		http.Error(w, "Failed to process Bend Works", http.StatusInternalServerError)
		return
	}

	taskIDPipeCutting, err := processPipeCutting(fileName, smartProcessID, engineerID)
	if err != nil {
		log.Printf("Error processing Pipe Cutting: %v\n", err)
		http.Error(w, "Failed to process Pipe Cutting", http.StatusInternalServerError)
		return
	}

	taskIDProducts, err := processProducts(fileName, smartProcessID, engineerID)
	if err != nil {
		log.Printf("Error processing products: %v\n", err)
		http.Error(w, "Failed to process products", http.StatusInternalServerError)
		return
	}
	arrayOfTasksIDs = append(arrayOfTasksIDs, taskIDLaserWorks, taskIDProducts, taskIDBendWorks, taskIDPipeCutting)

	// Проверяем наличие заполненных ячеек в столбце "Нанесение покрытий"
	if checkCoatingColumn(fileName) {
		err = pullCustomFieldInSmartProcess(1046, smartProcessID, "ufCrm6_1733264270", "да", arrayOfTasksIDs)
		if err != nil {
			log.Printf("Error updating smart process: %v\n", err)
			http.Error(w, "Failed to update smart process", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File processed successfully"))

}

// processLaserWorks обрабатывает столбец "Лазерные работы"
func processLaserWorks(fileName string, smartProcessID, engineerID int) (int, error) {
	return processTask(fileName, smartProcessID, engineerID, "Лазерные работы", 1)
}

// processBendWorks обрабатывает столбец "Гибочные работы"
func processBendWorks(fileName string, smartProcessID, engineerID int) (int, error) {
	return processTask(fileName, smartProcessID, engineerID, "Гибочные работы", 10)
}

// processPipeCutting обрабатывает столбец "Труборез"
func processPipeCutting(fileName string, smartProcessID, engineerID int) (int, error) {
	return processTask(fileName, smartProcessID, engineerID, "Труборез", 11)
}

// processTask универсальная функция для обработки задач
func processTask(fileName string, smartProcessID, engineerID int, taskType string, groupID int) (int, error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return 0, fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Реестр")
	if err != nil {
		return 0, fmt.Errorf("error reading rows: %v", err)
	}

	// Определяем индексы заголовков
	headers := map[string]int{
		"№ заказа":             -1,
		"Заказчик":             -1,
		"Менеджер":             -1,
		"Количество материала": -1,
		taskType:             -1,
		"Нанесение покрытий": -1,
		"Комментарий":        -1,
	}

	// Поиск заголовков
	for i, cell := range rows[0] {
		for header := range headers {
			if cell == header {
				headers[header] = i
				break
			}
		}
	}

	// Проверяем наличие всех необходимых заголовков
	for header, index := range headers {
		if index == -1 {
			return 0, fmt.Errorf("missing required header: %s", header)
		}
	}

	// Определяем конец таблицы
	var taskID int
	for _, row := range rows[1:] {
		// Если строка пуста, завершение обработки
		isEmptyRow := true
		for _, cell := range row {
			if cell != "" {
				isEmptyRow = false
				break
			}
		}
		if isEmptyRow {
			break
		}

		// Проверяем столбец
		if headers[taskType] >= len(row) || row[headers[taskType]] == "" {
			continue
		}

		// Создаём задачу, если ещё не создана
		if taskID == 0 {
			taskID, err = AddTaskToGroup(taskType, engineerID, groupID, 1046, smartProcessID)
			if err != nil {
				return 0, fmt.Errorf("error creating %s task: %v", taskType, err)
			}
		}

		// Создаём подзадачи
		customFields := CustomTaskFields{
			OrderNumber: row[headers["№ заказа"]],
			Customer:    row[headers["Заказчик"]],
			Manager:     row[headers["Менеджер"]],
			Quantity:    row[headers["Количество материала"]],
			Comment:     row[headers["Комментарий"]],
			Material:    row[headers[taskType]],
		}

		subTaskTitle := fmt.Sprintf("%s подзадача: %s", taskType, row[headers[taskType]])
		_, err := AddTaskToParentId(subTaskTitle, engineerID, groupID, taskID, customFields)
		if err != nil {
			log.Printf("Error creating %s subtask: %v\n", taskType, err)
			continue
		}
	}

	return taskID, nil
}

func HandlerProcessProducts(w http.ResponseWriter, r *http.Request) {
	products, err := processProducts("file.xlsx", 622, 149)
	if err != nil {
		log.Printf("Error processing products: %v\n", err)
	}
	fmt.Fprintf(w, "Products processed successfully: %v", products)
}

// processProducts обрабатывает столбцы "Производство" и "Нанесение покрытий"
func processProducts(fileName string, smartProcessID, engineerID int) (int, error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return 0, fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Реестр")
	if err != nil {
		return 0, fmt.Errorf("error reading rows: %v", err)
	}

	// Определяем индекс заголовка "Производство"
	headerIndex := -1
	for i, cell := range rows[0] {
		if cell == "Производство" {
			headerIndex = i
			break
		}
	}

	if headerIndex == -1 {
		return 0, fmt.Errorf("missing required header: Производство")
	}

	// Создаем основную задачу "Производство"
	taskID, err := AddTaskToGroup("Производство", engineerID, 2, 1046, smartProcessID)
	if err != nil {
		return 0, fmt.Errorf("error creating main production task: %v", err)
	}

	// Используем map для хранения уникальных строк
	uniqueRows := make(map[string]struct{})

	// Обрабатываем строки файла
	for _, row := range rows[1:] {
		// Проверяем, пуста ли строка
		if headerIndex >= len(row) || row[headerIndex] == "" {
			continue
		}

		productionCell := strings.TrimSpace(row[headerIndex]) // Текущая строка из "Производство"
		if _, exists := uniqueRows[productionCell]; !exists {
			uniqueRows[productionCell] = struct{}{}
		}
	}

	// Добавляем уникальные строки в чек-лист задачи
	for production := range uniqueRows {
		_, err := AddCheckListToTheTask(taskID, production)
		if err != nil {
			log.Printf("Error adding unique checklist item: %v\n", err)
		}
	}

	return taskID, nil
}

func parseCoatingCell(cell string) []string {
	// Аналогичная обработка для столбца "Нанесение покрытий"
	parts := strings.Split(cell, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}

// parseProductionCell парсит значение из столбца "Производство"
func parseProductionCell(cellValue string) []string {
	words := strings.Fields(cellValue)
	var buffer string

	// Для хранения уникальных элементов
	uniqueItems := make(map[string]struct{})
	var checklistItems []string

	for i, word := range words {
		// Начало нового элемента
		if strings.ToUpper(string(word[0])) == string(word[0]) {
			if buffer != "" {
				// Проверяем уникальность перед добавлением
				if _, exists := uniqueItems[buffer]; !exists {
					checklistItems = append(checklistItems, buffer)
					uniqueItems[buffer] = struct{}{}
				}
			}
			buffer = word
		} else {
			// Продолжение текущего элемента
			buffer += " " + word
		}

		// Последний элемент строки
		if i == len(words)-1 && buffer != "" {
			if _, exists := uniqueItems[buffer]; !exists {
				checklistItems = append(checklistItems, buffer)
				uniqueItems[buffer] = struct{}{}
			}
		}
	}

	return checklistItems
}

func mergeUniqueItems(items1, items2 []string) []string {
	uniqueSet := make(map[string]struct{})
	var mergedItems []string

	for _, item := range items1 {
		if _, exists := uniqueSet[item]; !exists {
			mergedItems = append(mergedItems, item)
			uniqueSet[item] = struct{}{}
		}
	}

	for _, item := range items2 {
		if _, exists := uniqueSet[item]; !exists {
			mergedItems = append(mergedItems, item)
			uniqueSet[item] = struct{}{}
		}
	}

	return mergedItems
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

// Добавить после downloadFile
