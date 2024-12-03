package laserflex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func LaserflexGetFile(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection is starting...")

	// Извлекаем параметры из URL
	queryParams := r.URL.Query()
	fileID := queryParams.Get("file_id")
	smartProcessIDStr := queryParams.Get("smartProcessID")
	engineerID := queryParams.Get("engineer_id")

	log.Printf("Engineer ID: %s", engineerID)

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

	// Парсим Excel файл
	taskIDs := map[string]int{} // Для хранения ID задач
	excelData, err := parseExcelFile(fileName)
	if err != nil {
		log.Printf("Error parsing Excel file: %v\n", err)
		http.Error(w, "Failed to parse Excel file", http.StatusInternalServerError)
		return
	}

	// Создаем задачи для каждой группы
	for taskType, taskData := range excelData {
		if len(taskData.Rows) > 0 { // Если есть данные в соответствующем столбце
			taskID, err := AddTaskToGroup(taskType, 149, taskData.GroupID, smartProcessID, 458)
			if err != nil {
				log.Printf("Error creating task for %s: %v\n", taskType, err)
				continue
			}
			log.Printf("%s Task created with ID: %d\n", taskType, taskID)
			taskIDs[taskType] = taskID

			// Создаем подзадачи
			for _, row := range taskData.Rows {
				customFields := CustomTaskFields{
					OrderNumber: row["№ заказа"],
					Customer:    row["Заказчик"],
					Manager:     row["Менеджер"],
					Quantity:    row["Количество материала"],
					Comment:     row["Комментарий"],
					Coating:     row["Нанесение покрытий"],
					Material:    row[taskType],
				}
				subTaskTitle := fmt.Sprintf("%s подзадача: %s", taskType, row[taskType])
				_, err := AddTaskToParentId(subTaskTitle, 149, taskData.GroupID, taskID, customFields)
				if err != nil {
					log.Printf("Error creating subtask for %s: %v\n", taskType, err)
					continue
				}
			}
		}
	}

	log.Println("All tasks and subtasks processed successfully")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File processed, tasks and subtasks added successfully"))
}

// Структура для хранения данных из Excel
type TaskData struct {
	GroupID int                 // ID группы (лазерные, труборез и т.д.)
	Rows    []map[string]string // Данные строк из Excel
}

func parseExcelFile(fileName string) (map[string]*TaskData, error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Реестр")
	if err != nil {
		return nil, fmt.Errorf("error reading rows: %v", err)
	}

	// Определяем индексы заголовков
	headers := map[string]int{
		"№ заказа":             -1,
		"Заказчик":             -1,
		"Менеджер":             -1,
		"Количество материала": -1,
		"Лазерные работы":      -1,
		"Труборез":             -1,
		"Гибочные работы":      -1,
		"Производство":         -1,
		"Нанесение покрытий":   -1,
		"Комментарий":          -1,
	}

	for i, cell := range rows[0] {
		for header := range headers {
			if strings.Contains(cell, header) {
				headers[header] = i
				break
			}
		}
	}

	// Проверяем наличие всех заголовков
	for header, index := range headers {
		if index == -1 {
			return nil, fmt.Errorf("missing required header: %s", header)
		}
	}

	// Инициализируем данные для задач
	excelData := map[string]*TaskData{
		"Лазерные работы": {GroupID: 1, Rows: []map[string]string{}},
		"Труборез":        {GroupID: 11, Rows: []map[string]string{}},
		"Гибочные работы": {GroupID: 10, Rows: []map[string]string{}},
		"Производство":    {GroupID: 2, Rows: []map[string]string{}},
	}

	// Парсим строки
	for _, cells := range rows[1:] {
		rowData := make(map[string]string)
		for header, index := range headers {
			if index >= 0 && index < len(cells) {
				rowData[header] = cells[index]
			}
		}

		// Проверяем заполненность данных и добавляем в соответствующие группы
		for taskType, task := range excelData {
			if rowData[taskType] != "" { // Если столбец не пустой
				task.Rows = append(task.Rows, rowData)
			}
		}
	}

	return excelData, nil
}

// вставить после downloadFile

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
