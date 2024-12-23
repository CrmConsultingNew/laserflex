package laserflex

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"strconv"
)

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
		"№ заказа":             -1,
		"Заказчик":             -1,
		"Количество материала": -1,
		taskType: -1,
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
			OrderNumber: row[headers["№ заказа"]],
			Customer:    row[headers["Заказчик"]],
			Quantity:    row[headers["Количество материала"]],
			Material:    row[headers[taskType]],
		}

		// Создаём задачу
		taskID, err := AddCustomTaskToParentId(taskTitle, 149, groupID, customFields, smartProcessID)
		if err != nil {
			log.Printf("Error creating %s task: %v\n", taskType, err)
			continue
		}
		taskIDs = append(taskIDs, taskID)
	}

	return taskIDs, nil
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
