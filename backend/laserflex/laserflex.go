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

// Структуры для задачи
type Task struct {
	Title         string   `json:"TITLE"`
	ResponsibleID int      `json:"RESPONSIBLE_ID"`
	GroupID       int      `json:"GROUP_ID,omitempty"`
	UfCrmTask     []string `json:"UF_CRM_TASK,omitempty"`
}

// Структура для создания задачи с полями
type TaskWithParent struct {
	Title         string           `json:"TITLE"`
	ResponsibleID int              `json:"RESPONSIBLE_ID"`
	GroupID       int              `json:"GROUP_ID,omitempty"`
	ParentID      int              `json:"PARENT_ID"`
	CustomFields  CustomTaskFields `json:"custom_fields,omitempty"`
}

// Структура для общего тела запроса
type TaskRequest struct {
	Fields map[string]interface{} `json:"fields"`
}

type TaskResponse struct {
	Result struct {
		Task struct {
			ID json.RawMessage `json:"id"`
		} `json:"task"`
	} `json:"result"`
}

// Структура для пользовательских полей задачи
type CustomTaskFields struct {
	Quantity          string `json:"UF_AUTO_552243496167,omitempty"` // Кол-во
	TemporaryOrderSum string `json:"UF_AUTO_555642596740,omitempty"` // Временная сумма заказа
	OrderNumber       string `json:"UF_AUTO_303168834495,omitempty"` // № заказа
	Customer          string `json:"UF_AUTO_876283676967,omitempty"` // Заказчик
	Manager           string `json:"UF_AUTO_794809224848,omitempty"` // Менеджер
	Material          string `json:"UF_AUTO_468857876599,omitempty"` // Материал
	Bend              string `json:"UF_AUTO_726724682983,omitempty"` // Гибка
	ProductionTask    string `json:"UF_AUTO_433735177517,omitempty"` // Произв. Задача
	Comment           string `json:"UF_AUTO_497907774817,omitempty"` // Комментарий
	Coating           string `json:"UF_AUTO_512869473370,omitempty"` // Покрытие
	AllowTimeTracking string `json:"ALLOW_TIME_TRACKING,omitempty"`
	TimeEstimate      int    `json:"TIME_ESTIMATE,omitempty"` //
}

func LaserflexGetFile(w http.ResponseWriter, r *http.Request) {
	//order_number={{№ заказа}}&deadline={{Срок сдачи}}
	log.Println("Connection is starting...")

	// Извлекаем параметры из URL
	queryParams := r.URL.Query()
	fileID := queryParams.Get("file_id")
	smartProcessIDStr := queryParams.Get("smartProcessID")
	orderNumber := queryParams.Get("order_number")

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

	// 2

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
		"Заказчик":             -1,
		"Количество материала": -1,
		taskType:               -1,
		"Время лазерных работ": -1, // Добавляем столбец для времени
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
		switch taskType {
		case "Лазерные работы":
			taskTitle = fmt.Sprintf("%s %s",
				orderNumber,
				row[headers[taskType]])
		case "Труборез":
			taskTitle = fmt.Sprintf("%s %s",
				orderNumber,
				row[headers[taskType]])
		case "Гибочные работы":
			taskTitle = fmt.Sprintf("Гибка %s %s",
				orderNumber,
				row[headers[taskType]])
		default:
			taskTitle = fmt.Sprintf("%s задача: %s",
				taskType, row[headers[taskType]])
		}

		customFields := CustomTaskFields{
			OrderNumber:       row[headers["№ заказа"]],
			Customer:          row[headers["Заказчик"]],
			Quantity:          row[headers["Количество материала"]],
			Material:          row[headers[taskType]],
			AllowTimeTracking: "Y",
			TimeEstimate:      timeEstimate, // Используем преобразованное значение
		}

		// Создаём задачу
		taskID, err := AddCustomTaskToParentId(orderNumber, taskTitle, 149, groupID, customFields, smartProcessID)
		if err != nil {
			log.Printf("Error creating %s task: %v\n", taskType, err)
			continue
		}
		taskIDs = append(taskIDs, taskID)
	}

	return taskIDs, nil
}

func AddCustomTaskToParentId(orderNumber string, title string, responsibleID, groupID int, customFields CustomTaskFields, elementID int) (int, error) {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "tasks.task.add"
	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	// Генерация ссылки смарт-процесса
	smartProcessLink := GenerateSmartProcessLink(1046, elementID)

	// Вычисление DEADLINE: Текущая дата + 13 часов
	currentTime := time.Now().Add(16 * time.Hour)
	deadline := currentTime.Format("02.01.2006T15:04:05")

	// Подготовка тела запроса
	requestBody := map[string]interface{}{
		"fields": map[string]interface{}{
			"TITLE":                title,
			"RESPONSIBLE_ID":       responsibleID,
			"GROUP_ID":             groupID,
			"UF_CRM_TASK":          []string{smartProcessLink},
			"UF_AUTO_303168834495": orderNumber,                     // № заказа
			"UF_AUTO_876283676967": []string{customFields.Customer}, // Заказчик
			"UF_AUTO_794809224848": []string{customFields.Manager},  // Менеджер
			"UF_AUTO_468857876599": []string{customFields.Material}, // Материал
			"UF_AUTO_497907774817": []string{customFields.Comment},  // Комментарий
			"UF_AUTO_552243496167": []string{customFields.Quantity}, // Кол-во
			"DEADLINE":             deadline,                        // DEADLINE: текущая дата + 13 часов
			"ALLOW_TIME_TRACKING":  "Y",                             // Обязательно строка "Y"
			"TIME_ESTIMATE":        customFields.TimeEstimate,       // Передаем как число
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
	var response TaskResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return 0, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Обрабатываем ID задачи
	var taskID int
	if err := json.Unmarshal(response.Result.Task.ID, &taskID); err != nil {
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
// Вставить после orderNumber

//deadline := queryParams.Get("deadline")

/*dealID := queryParams.Get("deal_id")
assignedByIdStr := queryParams.Get("assigned")

assignedById, err := strconv.Atoi(assignedByIdStr)
if err != nil {
	log.Printf("Error converting assigned ID to int: %v\n", err)
	http.Error(w, "Invalid assigned parameter", http.StatusBadRequest)
	return
}*/

// 2
// Вставить после downloadFile

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
