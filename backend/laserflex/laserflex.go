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
	//order_number={{№ заказа}}&deadline={{Срок сдачи}}
	log.Println("Connection is starting...")

	// Извлекаем параметры из URL
	queryParams := r.URL.Query()
	fileID := queryParams.Get("file_id")
	smartProcessIDStr := queryParams.Get("smartProcessID")
	orderNumber := queryParams.Get("order_number")
	//deadline := queryParams.Get("deadline")

	/*dealID := queryParams.Get("deal_id")
	assignedByIdStr := queryParams.Get("assigned")

	assignedById, err := strconv.Atoi(assignedByIdStr)
	if err != nil {
		log.Printf("Error converting assigned ID to int: %v\n", err)
		http.Error(w, "Invalid assigned parameter", http.StatusBadRequest)
		return
	}*/

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

	/*	// Чтение и обработка продуктов
		products, err := ReadXlsProductRows(fileName)
		if err != nil {
			log.Println("Error reading Excel file:", err)
			http.Error(w, "Failed to process Excel file", http.StatusInternalServerError)
			return
		}

		var productIDs []int
		var totalProductsPrice float64

		for _, product := range products {
			productID, err := AddProductsWithImage(product, "52")
			if err != nil {
				log.Printf("Error adding product %s: %v", product.Name, err)
				continue
			}
			productIDs = append(productIDs, productID)
			totalProductsPrice += product.Price * product.Quantity
		}

		var quantities, prices []float64
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

		for i, productId := range productIDs {
			quantity := quantities[i]
			err := AddCatalogDocumentElement(docId, productId, quantity)
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
	*/
	// Массив для всех ID задач

	var arrayOfTasksIDsLaser []int
	var arrayOfTasksIDsBend []int
	var arrayOfTasksIDsPipeCutting []int
	var arrayOfTasksIDsProducts []int
	// Обрабатываем задачи и собираем их ID
	if taskIDs, err := processLaserWorks(orderNumber, fileName, smartProcessID); err == nil {
		arrayOfTasksIDsLaser = append(arrayOfTasksIDsLaser, taskIDs...)
		log.Printf("ATTENT:!!!!!: arrayOfTasksIDsLaser ::: %v", arrayOfTasksIDsLaser)
	}

	if taskIDs, err := processBendWorks(orderNumber, fileName, smartProcessID); err == nil {
		arrayOfTasksIDsBend = append(arrayOfTasksIDsBend, taskIDs...)
		log.Printf("ATTENT:!!!!!: arrayOfTasksIDsBend ::: %v", arrayOfTasksIDsBend)
	}

	if taskIDs, err := processPipeCutting(orderNumber, fileName, smartProcessID); err == nil {
		arrayOfTasksIDsPipeCutting = append(arrayOfTasksIDsPipeCutting, taskIDs...)
		log.Printf("ATTENT:!!!!!: arrayOfTasksIDsPipeCutting ::: %v", arrayOfTasksIDsPipeCutting)
	}

	if taskIDs, err := processProducts(fileName, smartProcessID, 149); err == nil {
		arrayOfTasksIDsProducts = append(arrayOfTasksIDsProducts, taskIDs)
		log.Printf("ATTENT:!!!!!: arrayOfTasksIDsProducts ::: %v", arrayOfTasksIDsProducts)
	}

	// Лазерные работы ID
	err = pullCustomFieldInSmartProcess(false, 1046, smartProcessID, "ufCrm6_1734471089453", "да", arrayOfTasksIDsLaser)
	if err != nil {
		log.Printf("Error updating smart process: %v\n", err)
		http.Error(w, "Failed to update smart process", http.StatusInternalServerError)
		return
	}

	// Гибочные работы ID
	err = pullCustomFieldInSmartProcess(false, 1046, smartProcessID, "ufCrm6_1733265874338", "да", arrayOfTasksIDsBend) // Используем правильную переменную!
	if err != nil {
		log.Printf("Error updating smart process: %v\n", err)
		http.Error(w, "Failed to update smart process", http.StatusInternalServerError)
		return
	}

	// Труборез ID
	err = pullCustomFieldInSmartProcess(false, 1046, smartProcessID, "ufCrm6_1734471206084", "да", arrayOfTasksIDsPipeCutting) // Используем правильную переменную!
	if err != nil {
		log.Printf("Error updating smart process: %v\n", err)
		http.Error(w, "Failed to update smart process", http.StatusInternalServerError)
		return
	}

	// Проверяем наличие заполненных ячеек в столбце "Нанесение покрытий"

	// Проверяем наличие данных в столбце "Нанесение покрытий"
	if checkCoatingColumn(fileName) {
		// Если есть данные, получаем цвета из "Цвет/цинк"
		colors := parseSheetForColorColumn(fileName)
		_, err := AddTaskToGroupColor("Проверить наличие ЛКП на складе в ОМТС", 149, 12, 1046, smartProcessID, colors)
		if err != nil {
			log.Printf("Error creating task with colors: %v", err)
			http.Error(w, "Failed to create task with colors", http.StatusInternalServerError)
			return
		}

		// Обновляем смарт-процесс
		err = pullCustomFieldInSmartProcess(true, 1046, smartProcessID, "ufCrm6_1734478701624", "да", nil)
		if err != nil {
			log.Printf("Error updating smart process: %v\n", err)
			http.Error(w, "Failed to update smart process", http.StatusInternalServerError)
			return
		}
	} else {
		// Если данных нет, создаём задачу без цветов
		_, err := AddTaskToGroupColor("Задача в ОМТС с материалами из накладной", 149, 12, 1046, smartProcessID, nil)
		if err != nil {
			log.Printf("Error creating task without colors: %v", err)
			http.Error(w, "Failed to create task without colors", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File processed successfully"))
}

// processLaserWorks обрабатывает столбец "Лазерные работы"
func processLaserWorks(orderNumber string, fileName string, smartProcessID int) ([]int, error) {
	return processTaskCustom(orderNumber, fileName, smartProcessID, "Лазерные работы", 1)
}

// processBendWorks обрабатывает столбец "Гибочные работы"
func processBendWorks(orderNumber string, fileName string, smartProcessID int) ([]int, error) {
	return processTaskCustom(orderNumber, fileName, smartProcessID, "Гибочные работы", 10)
}

// processPipeCutting обрабатывает столбец "Труборез"
func processPipeCutting(orderNumber string, fileName string, smartProcessID int) ([]int, error) {
	return processTaskCustom(orderNumber, fileName, smartProcessID, "Труборез", 11)
}

// processTaskCustom использует AddCustomTaskToParentId для обработки задач
func processTaskCustom(orderNumber string, fileName string, smartProcessID int, taskType string, groupID int) ([]int, error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Реестр")
	if err != nil {
		return nil, fmt.Errorf("error reading rows: %v", err)
	}

	headers := map[string]int{
		"Гибочные работы":      -1,
		"Лазерные работы":      -1,
		"№ заказа":             -1,
		"Заказчик":             -1,
		"Количество материала": -1,
		"Время лазерных работ": -1,
		"Время гибочных работ": -1,
		"Труборез":             -1,
		"Время труборез":       -1,
		taskType:               -1,
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
			return nil, fmt.Errorf("missing required header: %s", header)
		}
	}

	// Массив для хранения ID созданных задач
	var taskIDs []int

	// Обработка строк
	for _, row := range rows[1:] {
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

		// Пропускаем строки без данных
		if headers[taskType] >= len(row) || row[headers[taskType]] == "" {
			continue
		}

		var taskTitle string
		var timeEstimate int

		switch taskType {
		case "Лазерные работы":
			taskTitle = fmt.Sprintf("%s %s",
				orderNumber,
				row[headers["Лазерные работы"]])
			timeEstimate = convertToSeconds(row[headers["Время лазерных работ"]])
		case "Труборез":
			taskTitle = fmt.Sprintf("%s %s",
				orderNumber,
				row[headers["Труборез"]])
			timeEstimate = convertToSeconds(row[headers["Время труборез"]])
		case "Гибочные работы":
			taskTitle = fmt.Sprintf("Гибка %s %s",
				orderNumber,
				row[headers["Лазерные работы"]])
			timeEstimate = convertToSeconds(row[headers["Время гибочных работ"]])
		default:
			taskTitle = fmt.Sprintf("%s задача: %s",
				taskType, row[headers[taskType]])
			timeEstimate = 0 // Для других типов время по умолчанию
		}

		customFields := CustomTaskFields{
			OrderNumber: row[headers["№ заказа"]],
			Customer:    row[headers["Заказчик"]],
			Quantity:    row[headers["Количество материала"]],
			Material:    row[headers[taskType]],
		}

		// Создаём задачу
		taskID, err := AddCustomTaskToParentId(taskTitle, 149, groupID, customFields, smartProcessID, timeEstimate)
		if err != nil {
			log.Printf("Error creating %s task: %v\n", taskType, err)
			continue
		}
		taskIDs = append(taskIDs, taskID)
	}

	return taskIDs, nil
}

func convertToSeconds(cellValue string) int {
	// Преобразование времени в минуты в секунды
	value, err := strconv.Atoi(cellValue)
	if err != nil {
		return 0
	}
	return value * 60
}

func HandlerProcessProducts(w http.ResponseWriter, r *http.Request) {
	products, err := processProducts("file.xlsx", 688, 149)
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

	// Определяем индексы заголовков
	headers := map[string]int{
		"Производство":         -1,
		"Нанесение покрытий":   -1,
		"№ заказа":             -1,
		"Заказчик":             -1,
		"Менеджер":             -1,
		"Комментарий":          -1,
		"Количество материала": -1,
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

	// ID основной задачи "Производство"
	taskID, err := AddTaskToGroup("Производство", engineerID, 2, 1046, smartProcessID)
	if err != nil {
		return 0, fmt.Errorf("error creating main production task: %v", err)
	}

	// Используем map для проверки уникальности
	uniqueChecklistItems := make(map[string]struct{})

	// Обработка строк и добавление чек-листов
	for _, row := range rows[1:] {
		// Проверяем пустоту строки
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

		// Получаем значения ячеек
		productionCell := row[headers["Производство"]]
		coatingCell := row[headers["Нанесение покрытий"]]

		// Проверяем и добавляем элементы из "Производство"
		if productionCell != "" {
			if _, exists := uniqueChecklistItems[productionCell]; !exists {
				uniqueChecklistItems[productionCell] = struct{}{}
				_, err := AddCheckListToTheTask(taskID, productionCell)
				if err != nil {
					log.Printf("Error adding checklist item from 'Производство': %v\n", err)
				}
			}
		}

		// Проверяем и добавляем элементы из "Нанесение покрытий"
		if coatingCell != "" {
			if _, exists := uniqueChecklistItems[coatingCell]; !exists {
				uniqueChecklistItems[coatingCell] = struct{}{}
				_, err := AddCheckListToTheTask(taskID, coatingCell)
				if err != nil {
					log.Printf("Error adding checklist item from 'Нанесение покрытий': %v\n", err)
				}
			}
		}
	}

	return taskID, nil
}

func processCoatingTasks(fileName string, smartProcessID, engineerID int) error {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Реестр")
	if err != nil {
		return fmt.Errorf("error reading rows: %v", err)
	}

	// Определяем индексы заголовков
	headers := map[string]int{
		"Нанесение покрытий": -1,
		"№ заказа":           -1,
		"Заказчик":           -1,
		"Цвет / Цинк":        -1,
	}

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
			return fmt.Errorf("missing required header: %s", header)
		}
	}

	// Уникальные значения для "Цвет / Цинк"
	uniqueColors := make(map[string]struct{})

	// Обработка строк
	for _, row := range rows[1:] {
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

		orderNumber := row[headers["№ заказа"]]
		customer := row[headers["Заказчик"]]
		coating := row[headers["Нанесение покрытий"]]
		color := row[headers["Цвет / Цинк"]]

		if coating != "" {
			if _, exists := uniqueColors[color]; !exists {
				uniqueColors[color] = struct{}{}
			}

			taskTitle := fmt.Sprintf("Проверить наличие ЛКП на складе в ОМТС по Заказу %s", orderNumber)

			customFields := CustomTaskFields{
				OrderNumber: orderNumber,
				Customer:    customer,
			}

			// Создаём задачу, передавая массив уникальных цветов
			colorArray := make([]string, 0, len(uniqueColors))
			for color := range uniqueColors {
				colorArray = append(colorArray, color)
			}

			_, err := AddCustomCoatingTask(taskTitle, engineerID, 12, customFields, smartProcessID, colorArray)
			if err != nil {
				log.Printf("Error creating coating task: %v\n", err)
				continue
			}
		}
	}

	log.Printf("Processing coating tasks completed successfully")
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
			taskID, err = AddTaskToGroup(taskType, 149, groupID, 1046, smartProcessID)
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
		_, err := AddTaskToParentId(subTaskTitle, 149, groupID, taskID, customFields)
		if err != nil {
			log.Printf("Error creating %s subtask: %v\n", taskType, err)
			continue
		}
	}

	return taskID, nil
}

func parseCoatingCell(cellValue string) []string {
	words := strings.Fields(cellValue)
	var checklistItems []string
	var buffer string

	for i, word := range words {
		if strings.ToUpper(string(word[0])) == string(word[0]) {
			if buffer != "" {
				checklistItems = append(checklistItems, buffer)
			}
			buffer = word
		} else {
			buffer += " " + word
		}

		if i == len(words)-1 && buffer != "" {
			checklistItems = append(checklistItems, buffer)
		}
	}

	return checklistItems
}

// parseProductionCell парсит значение из столбца "Производство"
func parseProductionCell(cellValue string) []string {
	words := strings.Fields(cellValue)
	var checklistItems []string
	var buffer string

	for i, word := range words {
		if strings.ToUpper(string(word[0])) == string(word[0]) {
			if buffer != "" {
				checklistItems = append(checklistItems, buffer)
			}
			buffer = word
		} else {
			buffer += " " + word
		}

		if i == len(words)-1 && buffer != "" {
			checklistItems = append(checklistItems, buffer)
		}
	}

	return checklistItems
}
