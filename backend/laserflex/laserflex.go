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
	"time"
)

func LaserflexGetFile(w http.ResponseWriter, r *http.Request) {
	//order_number={{№ заказа}}&deadline={{Срок сдачи}}
	log.Println("Connection is starting...")

	// Извлекаем параметры из URL
	queryParams := r.URL.Query()
	fileID := queryParams.Get("file_id")
	smartProcessIDStr := queryParams.Get("smartProcessID")
	orderNumber := queryParams.Get("order_number")
	engineerStr := queryParams.Get("engineer_id")
	engineerId, _ := strconv.Atoi(engineerStr)
	deadline := queryParams.Get("deadline")
	// 1
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

	//2

	var arrayOfTasksIDsLaser []int
	var arrayOfTasksIDsBend []int
	var arrayOfTasksIDsPipeCutting []int
	var arrayOfTasksIDsProducts []int
	// Обрабатываем задачи и собираем их ID
	if taskIDs, err := processLaserWorks(orderNumber, fileName, smartProcessID, deadline); err == nil {
		arrayOfTasksIDsLaser = append(arrayOfTasksIDsLaser, taskIDs...)
	}

	if taskIDs, err := processBendWorks(orderNumber, fileName, smartProcessID, deadline); err == nil {
		arrayOfTasksIDsBend = append(arrayOfTasksIDsBend, taskIDs...)
	}

	if taskIDs, err := processPipeCutting(orderNumber, fileName, smartProcessID, deadline); err == nil {
		arrayOfTasksIDsPipeCutting = append(arrayOfTasksIDsPipeCutting, taskIDs...)
	}

	if taskIDs, err := processProducts(fileName, smartProcessID, engineerId); err == nil {
		arrayOfTasksIDsProducts = append(arrayOfTasksIDsProducts, taskIDs)
	}

	log.Printf("Laser task IDs: %v", arrayOfTasksIDsLaser)
	log.Printf("Bend task IDs: %v", arrayOfTasksIDsBend)
	log.Printf("Pipe Cutting task IDs: %v", arrayOfTasksIDsPipeCutting)
	log.Printf("Products task IDs: %v", arrayOfTasksIDsProducts)

	if len(arrayOfTasksIDsLaser) > 0 {
		err = pullCustomFieldInSmartProcess(false, 1046, smartProcessID, "ufCrm6_1734471089453", "да", arrayOfTasksIDsLaser)
		if err != nil {
			log.Printf("Error updating smart process for Laser tasks: %v\n", err)
		}
	}

	if len(arrayOfTasksIDsBend) > 0 {
		err = pullCustomFieldInSmartProcess(false, 1046, smartProcessID, "ufCrm6_1733265874338", "да", arrayOfTasksIDsBend)
		if err != nil {
			log.Printf("Error updating smart process for Bend tasks: %v\n", err)
		}
	}

	if len(arrayOfTasksIDsPipeCutting) > 0 {
		err = pullCustomFieldInSmartProcess(false, 1046, smartProcessID, "ufCrm6_1734471206084", "да", arrayOfTasksIDsPipeCutting)
		if err != nil {
			log.Printf("Error updating smart process for Pipe Cutting tasks: %v\n", err)
		}
	}

	// Проверяем наличие данных в столбце "Нанесение покрытий"
	if CheckCoatingColumn(fileName) {
		// Если есть данные, получаем цвета из "Цвет/цинк"
		colors := ParseSheetForColorColumn(fileName)
		_, err := AddTaskToGroupColor(orderNumber, "Проверить наличие ЛКП на складе в ОМТС", engineerId, 12, 1046, smartProcessID, colors)
		if err != nil {
			log.Printf("Error creating task with colors: %v", err)
			http.Error(w, "Failed to create task with colors", http.StatusInternalServerError)
			return
		}

	} else {
		// Если данных нет, создаём задачу без цветов
		_, err := AddTaskToGroupColor(orderNumber, "Задача в ОМТС с материалами из накладной", engineerId, 12, 1046, smartProcessID, nil)
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
func processLaserWorks(orderNumber string, fileName string, smartProcessID int, deadline string) ([]int, error) {
	return processTaskCustom(orderNumber, fileName, smartProcessID, "Лазерные работы", 1, deadline)
}

// processBendWorks обрабатывает столбец "Гибочные работы"
func processBendWorks(orderNumber string, fileName string, smartProcessID int, deadline string) ([]int, error) {
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
		"Заказчик":             -1,
		"Количество материала": -1,
		"Гибочные работы":      -1,
		"Лазерные работы":      -1,
	}

	// Поиск заголовков
	for i, cell := range rows[0] {
		if _, ok := headers[cell]; ok {
			headers[cell] = i
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

		if headers["Гибочные работы"] >= len(row) || row[headers["Гибочные работы"]] == "" {
			continue
		}

		// Получаем значения для заголовка задачи и пользовательских полей
		laserWorksValue := ""
		if headers["Лазерные работы"] < len(row) {
			laserWorksValue = row[headers["Лазерные работы"]]
		}

		materialValue := ""
		if headers["Количество материала"] < len(row) {
			materialValue = row[headers["Лазерные работы"]]
		}

		timeEstimateStr := row[headers["Гибочные работы"]]
		timeEstimate, err := strconv.Atoi(timeEstimateStr)
		if err != nil {
			log.Printf("Error converting time estimate '%s' to int: %v", timeEstimateStr, err)
			continue
		}

		taskTitle := fmt.Sprintf("Гибка %s %s", orderNumber, laserWorksValue)

		customFields := CustomTaskFields{
			OrderNumber:       orderNumber,
			Customer:          row[headers["Заказчик"]],
			Material:          materialValue, // Используем значение из столбца "Количество материала"
			AllowTimeTracking: "Y",
			TimeEstimate:      timeEstimate,
		}

		// Создаём задачу
		taskID, err := AddCustomTaskToParentId(orderNumber, taskTitle, 149, 10, customFields, smartProcessID, deadline)
		if err != nil {
			log.Printf("Error creating BendWorks task: %v\n", err)
			continue
		}
		taskIDs = append(taskIDs, taskID)
	}

	return taskIDs, nil
}

// processPipeCutting обрабатывает столбец "Труборез"
func processPipeCutting(orderNumber string, fileName string, smartProcessID int, deadline string) ([]int, error) {
	return processTaskCustom(orderNumber, fileName, smartProcessID, "Труборез", 11, deadline)
}

// processTaskCustom использует AddCustomTaskToParentId для обработки задач
func processTaskCustom(orderNumber string, fileName string, smartProcessID int, taskType string, groupID int, deadline string) ([]int, error) {
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
		"Заказчик":             -1,
		"Количество материала": -1,
		taskType:               -1,
		"Лазерные работы":      -1, // Добавляем столбец для "Лазерные работы"
		"Время лазерных работ": -1,
	}

	// Поиск заголовков
	for i, cell := range rows[0] {
		if _, ok := headers[cell]; ok {
			headers[cell] = i
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

		if headers[taskType] >= len(row) || row[headers[taskType]] == "" {
			continue
		}

		// Преобразуем время из строки в int
		timeEstimateStr := row[headers["Время лазерных работ"]]
		timeEstimate, err := strconv.Atoi(timeEstimateStr)
		if err != nil {
			log.Printf("Error converting time estimate '%s' to int: %v", timeEstimateStr, err)
			continue
		}

		// Формируем заголовок задачи на основе taskType
		taskTitle := ""
		if taskType == "Гибочные работы" {
			cellOfLaserWorksValue := ""
			if headers["Лазерные работы"] < len(row) {
				cellOfLaserWorksValue = row[headers["Лазерные работы"]]
			}
			taskTitle = fmt.Sprintf("Гибка %s %s", orderNumber, cellOfLaserWorksValue)
		} else {
			taskTitle = fmt.Sprintf("%s %s", orderNumber, row[headers[taskType]])
		}

		customFields := CustomTaskFields{
			OrderNumber:       row[headers["Заказчик"]],
			Customer:          row[headers["Заказчик"]],
			Quantity:          row[headers["Количество материала"]],
			Material:          row[headers[taskType]],
			AllowTimeTracking: "Y",
			TimeEstimate:      timeEstimate, // Используем преобразованное значение
		}

		// Создаём задачу
		taskID, err := AddCustomTaskToParentId(orderNumber, taskTitle, 149, groupID, customFields, smartProcessID, deadline)
		if err != nil {
			log.Printf("Error creating %s task: %v\n", taskType, err)
			continue
		}
		taskIDs = append(taskIDs, taskID)
	}

	return taskIDs, nil
}

func AddCustomTaskToParentId(orderNumber string, title string, responsibleID, groupID int, customFields CustomTaskFields, elementID int, deadline string) (int, error) {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "tasks.task.add"
	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	// Парсим входящую дату
	parsedDate, err := time.Parse("02.01.2006", deadline)
	if err != nil {
		return 0, fmt.Errorf("invalid deadline format: %v", err)
	}

	// Добавляем 16 часов к дате
	deadlineWithTime := parsedDate.Add(16 * time.Hour).Format("02.01.2006T15:04:05")

	// Генерация ссылки смарт-процесса
	smartProcessLink := GenerateSmartProcessLink(1046, elementID)

	// Подготовка тела запроса
	requestBody := map[string]interface{}{
		"fields": map[string]interface{}{
			"TITLE":                title,
			"RESPONSIBLE_ID":       responsibleID,
			"GROUP_ID":             groupID,
			"UF_CRM_TASK":          []string{smartProcessLink},
			"UF_AUTO_303168834495": []string{orderNumber},
			"UF_AUTO_876283676967": []string{customFields.Customer},
			"UF_AUTO_794809224848": []string{customFields.Manager},
			"UF_AUTO_468857876599": []string{customFields.Material},
			"UF_AUTO_497907774817": []string{customFields.Comment},
			"UF_AUTO_552243496167": []string{customFields.Quantity},
			"DEADLINE":             deadlineWithTime,
			"ALLOW_TIME_TRACKING":  "Y",
			"TIME_ESTIMATE":        customFields.TimeEstimate * 60,
		},
	}

	// Сериализация тела запроса
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return 0, fmt.Errorf("error marshalling request body: %v", err)
	}

	log.Printf("Request Body: %s", string(jsonData))

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	log.Printf("Response from Bitrix24: %s\n", string(responseData))

	// Разбираем ответ
	var response struct {
		Result struct {
			Task struct {
				ID string `json:"id"`
			} `json:"task"`
		} `json:"result"`
	}
	if err := json.Unmarshal(responseData, &response); err != nil {
		return 0, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Преобразуем строковый ID в число
	taskID, err := strconv.Atoi(response.Result.Task.ID)
	if err != nil {
		return 0, fmt.Errorf("error parsing task id: %v", err)
	}

	return taskID, nil
}

func pullCustomFieldInSmartProcess(checkCoating bool, entityTypeId, smartProcessID int, fieldName, fieldValue string, tasksIDs []int) error {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "crm.item.update"
	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	if len(tasksIDs) == 0 {
		return fmt.Errorf("tasksIDs array is empty, cannot update smart process")
	}
	log.Printf("Updating smart process ID %d with tasksIDs: %v", smartProcessID, tasksIDs)

	// Преобразование tasksIDs в []string
	stringTasksIDs := make([]string, len(tasksIDs))
	for i, id := range tasksIDs {
		stringTasksIDs[i] = strconv.Itoa(id)
	}

	if checkCoating == true {

	}
	// Обновляем значение полей в запросе
	requestBody := map[string]interface{}{
		"entityTypeId": entityTypeId,
		"id":           smartProcessID,
		"fields": map[string]interface{}{
			fieldName: stringTasksIDs, // Используем динамическое имя поля
		},
	}
	if checkCoating == true {
		requestBody = map[string]interface{}{
			"entityTypeId": entityTypeId,
			"id":           smartProcessID,
			"fields": map[string]interface{}{
				"ufCrm6_1733264270": "да", // Используем динамическое имя поля
			},
		}
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %v", err)
	}
	log.Printf("Request Body: %s", string(jsonData))

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
	log.Printf("Response from pullCustomFieldInSmartProcess: %s", string(responseData))

	var response map[string]interface{}
	if err := json.Unmarshal(responseData, &response); err != nil {
		return fmt.Errorf("error unmarshalling response: %v", err)
	}

	if _, ok := response["error"]; ok {
		return fmt.Errorf("failed to update smart process: %s", response["error_description"])
	}

	log.Printf("Smart process updated successfully for ID: %d with tasks: %v", smartProcessID, tasksIDs)
	return nil
}

// 1

/*dealID := queryParams.Get("deal_id")
assignedByIdStr := queryParams.Get("assigned")

assignedById, err := strconv.Atoi(assignedByIdStr)
if err != nil {
log.Printf("Error converting assigned ID to int: %v\n", err)
http.Error(w, "Invalid assigned parameter", http.StatusBadRequest)
return
}*/

// 2

// Чтение и обработка продуктов
/*products, err := ReadXlsProductRows(fileName)
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
}*/
