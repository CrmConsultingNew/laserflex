package laserflex

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

/*func HandlerProcessProducts(w http.ResponseWriter, r *http.Request) {
	products, err := processProducts("file.xlsx", 688, 149)
	if err != nil {
		log.Printf("Error processing products: %v\n", err)
	}
	fmt.Fprintf(w, "Products processed successfully: %v", products)
}*/

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
func processTask(orderNumber, client, fileName string, smartProcessID, engineerID int, taskType string, groupID int) (int, error) {
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
			taskID, err = AddTaskToGroup(orderNumber, client, taskType, 149, groupID, 1046, smartProcessID)
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

type Product struct {
	Name        string  // Наименование (столбец A)
	Quantity    float64 // Количество (столбец B)
	Price       float64 // Цена (столбец C)
	ImageBase64 string  // Изображение в Base64 (столбец D)
	Material    float64 // Материал (столбец E)
	Laser       float64 // Лазер (столбец F)
	Bend        float64 // Гиб (столбец G)
	Weld        float64 // Свар (столбец H)
	Paint       float64 // Окраска (столбец I)
	Production  float64 // Производство (Сумма столбцов H, J, K, L, M, N, O)
	AddP        float64 // Допы П (столбец N)
	AddL        float64 // Допы Л (столбец O)
	PipeCutting float64 // Труборез (столбец P)
}

// Специальная функция для обработки столбца Price
func parsePrice(input string) float64 {
	fmt.Printf("Original Price input: '%s'\n", input)

	// Удаление пробелов между цифрами
	re := regexp.MustCompile(`(\d)\s+(\d)`)
	input = re.ReplaceAllString(input, "$1$2")
	fmt.Printf("After removing spaces: '%s'\n", input)

	// Замена всех запятых на точки
	input = strings.ReplaceAll(input, ",", ".")

	// Проверка на наличие более одной точки
	if strings.Count(input, ".") > 1 {
		// Оставляем только последнюю точку
		parts := strings.Split(input, ".")
		input = strings.Join(parts[:len(parts)-1], "") + "." + parts[len(parts)-1]
		fmt.Printf("After fixing dots: '%s'\n", input)
	}

	// Преобразование в float64
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Printf("Error parsing Price: %v\n", err)
		return 0
	}

	// Округляем до двух знаков после запятой
	value = math.Round(value*100) / 100
	fmt.Printf("Parsed and rounded Price: %f\n", value)

	return value
}

// Функция для обработки чисел в других столбцах
func parseFloatOrInt(input string) float64 {
	fmt.Printf("Original input: '%s'\n", input)

	// Убираем пробелы между цифрами
	re := regexp.MustCompile(`(\d)\s+(\d)`)
	input = re.ReplaceAllString(input, "$1$2")
	//fmt.Printf("After removing spaces: '%s'\n", input)

	// Заменяем запятую на точку
	input = strings.ReplaceAll(input, ",", ".")
	//fmt.Printf("After replacing commas: '%s'\n", input)

	// Пробуем преобразовать в float64
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Printf("Warning: unable to parse float or int from string '%s': %v\n", input, err)
		return 0
	}

	//fmt.Printf("Parsed float: %f\n", value)
	return value
}

// Функция для чтения строк из Excel
func ReadXlsProductRows(filename string) ([]Product, error) {
	fmt.Println("Processing Excel file...")

	// Открываем файл Excel
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	// Читаем строки из листа "Статистика"
	rows, err := f.GetRows("Статистика")
	if err != nil {
		return nil, fmt.Errorf("error reading rows: %v", err)
	}

	var products []Product

	// Обрабатываем каждую строку
	for i, cells := range rows {
		if i == 0 || len(cells) < 16 { // Пропускаем заголовок и проверяем минимальное количество столбцов
			continue
		}

		// Проверяем первую ячейку на условия завершения
		if len(cells) > 0 {
			name := strings.TrimSpace(cells[0]) // Убираем лишние пробелы
			if name == "" || strings.EqualFold(name, "Доставка") || strings.Contains(strings.ToLower(name), "общее") {
				fmt.Printf("Terminating parsing at row %d: Name='%s'\n", i+1, name)
				break
			}
		}

		// Получение Base64 строки изображения из ячейки
		imageBase64 := ""
		imageData, err := getImageBase64FromExcel(f, "Статистика", fmt.Sprintf("D%d", i+1))
		if err == nil {
			imageBase64 = imageData
		} else {
			fmt.Printf("Warning: unable to get image for row %d: %v\n", i+1, err)
		}

		// Суммируем для Production
		production := 0.0
		for _, colIndex := range []int{7, 9, 10, 11, 12, 13, 14} {
			if colIndex < len(cells) && cells[colIndex] != "" {
				production += parseFloatOrInt(cells[colIndex])
			} else {
				fmt.Printf("Warning: missing or empty value at column %d, row %d\n", colIndex, i+1)
			}
		}

		// Создаём объект Product
		product := Product{
			Name:        cells[0],
			Quantity:    parseFloatOrInt(cells[1]),
			Price:       parsePrice(cells[2]), // Используем parsePrice для обработки столбца C
			ImageBase64: imageBase64,
			Material:    parseFloatOrInt(cells[4]),
			Laser:       parseFloatOrInt(cells[5]),
			Bend:        parseFloatOrInt(cells[6]),
			Weld:        parseFloatOrInt(cells[7]),
			Paint:       parseFloatOrInt(cells[8]),
			Production:  production,
			AddP:        parseFloatOrInt(cells[13]),
			AddL:        parseFloatOrInt(cells[14]),
			PipeCutting: parseFloatOrInt(cells[15]),
		}

		products = append(products, product)
		fmt.Printf("Parsed Product: %+v\n", product)
	}

	fmt.Println("Excel processing completed.")
	return products, nil
}

// getImageBase64FromExcel извлекает изображение из ячейки и возвращает его в виде строки Base64
func getImageBase64FromExcel(f *excelize.File, sheet, cell string) (string, error) {
	pictures, err := f.GetPictures(sheet, cell)
	if err != nil {
		return "", fmt.Errorf("error extracting image from cell %s: %v", cell, err)
	}

	if len(pictures) == 0 {
		return "", fmt.Errorf("no images found in cell %s", cell)
	}

	// Извлекаем данные первого изображения
	imageData := pictures[0].File

	// Кодируем изображение в Base64
	return base64.StdEncoding.EncodeToString(imageData), nil
}

// Универсальная структура для хранения данных из строки
type ParsedData struct {
	LaserWorks  *WorkGroup `json:"laser_works,omitempty"`  // Для строк с Лазерными работами
	PipeCutting *WorkGroup `json:"pipe_cutting,omitempty"` // Для строк с Труборезом
	BendWorks   *WorkGroup `json:"bend_works,omitempty"`   // Для строк с Гибочными работами
	Production  *WorkGroup `json:"production,omitempty"`   // Для строк с Производством
}

// WorkGroup структура для хранения данных группы
type WorkGroup struct {
	GroupID int               `json:"group_id"`
	Data    map[string]string `json:"data"`
}

// Функция для чтения таблицы и разделения данных по условиям
func ReadXlsRegistryWithConditions(filename string) ([]ParsedData, error) {
	fmt.Println("Processing Registry Excel file...")

	f, err := excelize.OpenFile(filename)
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
		"Время лазерных работ": -1,
		"Производство":         -1,
		"Нанесение покрытий":   -1, // Добавляем столбец "Нанесение покрытий"
		"Комментарий":          -1,
	}

	// Найдем индексы всех необходимых заголовков
	for i, cell := range rows[0] {
		for header := range headers {
			if strings.Contains(cell, header) {
				headers[header] = i
				break
			}
		}
	}

	// Проверим, что все необходимые заголовки найдены
	for header, index := range headers {
		if index == -1 {
			return nil, fmt.Errorf("missing required header: %s", header)
		}
	}

	// Парсинг строк
	var parsedRows []ParsedData
	for _, cells := range rows[1:] {
		// Проверяем, если строка полностью пустая
		isEmptyRow := true
		for _, index := range headers {
			if index < len(cells) && cells[index] != "" {
				isEmptyRow = false
				break
			}
		}
		if isEmptyRow {
			break
		}

		// Создаем структуру для хранения данных
		parsedRow := ParsedData{}

		// Проверяем условия для каждого столбца
		if value := getValue(cells, headers["Лазерные работы"]); value != "" {
			parsedRow.LaserWorks = &WorkGroup{
				GroupID: 1,
				Data: extractData(cells, headers, []string{
					"№ заказа", "Заказчик", "Менеджер", "Количество материала", "Лазерные работы", "Нанесение покрытий", "Комментарий",
				}),
			}
		}

		if value := getValue(cells, headers["Труборез"]); value != "" {
			parsedRow.PipeCutting = &WorkGroup{
				GroupID: 11,
				Data: extractData(cells, headers, []string{
					"№ заказа", "Заказчик", "Менеджер", "Количество материала", "Труборез", "Нанесение покрытий", "Комментарий",
				}),
			}
		}

		if value := getValue(cells, headers["Гибочные работы"]); value != "" {
			parsedRow.BendWorks = &WorkGroup{
				GroupID: 10,
				Data: extractData(cells, headers, []string{
					"№ заказа", "Заказчик", "Менеджер", "Количество материала", "Гибочные работы", "Нанесение покрытий", "Комментарий",
				}),
			}
		}

		if value := getValue(cells, headers["Производство"]); value != "" {
			parsedRow.Production = &WorkGroup{
				GroupID: 2,
				Data: extractData(cells, headers, []string{
					"№ заказа", "Заказчик", "Менеджер", "Количество материала", "Производство", "Нанесение покрытий", "Комментарий",
				}),
			}
		}

		// Если данные для текущей строки соответствуют хотя бы одному условию, добавляем строку
		if parsedRow.LaserWorks != nil || parsedRow.PipeCutting != nil || parsedRow.BendWorks != nil || parsedRow.Production != nil {
			parsedRows = append(parsedRows, parsedRow)
		}
	}

	// Выводим результат
	//fmt.Printf("\nParsed Rows:\n%v\n", parsedRows)

	return parsedRows, nil
}

// Извлекает данные для указанных столбцов
func extractData(cells []string, headers map[string]int, columns []string) map[string]string {
	data := make(map[string]string)
	for _, column := range columns {
		index := headers[column]
		if index >= 0 && index < len(cells) {
			data[column] = cells[index]
		}
	}
	return data
}

// Получает значение из ячейки или возвращает пустую строку, если индекс выходит за пределы
func getValue(cells []string, index int) string {
	if index >= 0 && index < len(cells) {
		return cells[index]
	}
	return ""
}

func AddCatalogDocument(dealID string, assignedById int, totalProductsPrice float64) (int, error) {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "catalog.document.add"

	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	// Формируем тело запроса
	requestBody := map[string]interface{}{
		"fields": map[string]interface{}{
			"docType":       "S",
			"responsibleId": assignedById,
			"createdBy":     assignedById,
			"currency":      "RUB",
			"status":        "S",
			"statusBy":      assignedById,
			"total":         fmt.Sprintf("%.2f", totalProductsPrice), // Округляем до двух знаков
			"title":         dealID,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return 0, fmt.Errorf("error marshalling request body: %v", err)
	}

	// Создаем HTTP запрос
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
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

	// Парсим ответ
	var response struct {
		Result struct {
			Document struct {
				ID int `json:"id"`
			} `json:"document"`
		} `json:"result"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}

	if err := json.Unmarshal(responseData, &response); err != nil {
		return 0, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Проверяем наличие ошибок в ответе
	if response.Error != "" {
		return 0, fmt.Errorf("Ошибка: %s", response.ErrorDescription)
	}

	// Логируем успешное добавление документа
	log.Printf("Документ успешно добавлен: ID=%d\n", response.Result.Document.ID)

	// Возвращаем ID созданного документа
	return response.Result.Document.ID, nil
}

func AddCatalogDocumentElement(documentID int, productId int, quantity float64) error {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "catalog.document.element.add"

	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	requestBody := map[string]interface{}{
		"fields": map[string]interface{}{
			"docId":     documentID,
			"elementId": productId,
			"storeTo":   "1",           // ID склада
			"amount":    int(quantity), // Количество из массива
		},
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
		log.Printf("Full response on error: %v", string(responseData))
		return fmt.Errorf("ошибка: %s", response["error_description"])
	}

	if _, ok := response["result"]; !ok {
		log.Printf("Unexpected response structure: %v", string(responseData))
		return fmt.Errorf("missing 'result' in response")
	}

	log.Printf("Товар добавлен в документ: ID документа: %s, ID товара: %d\n", documentID, productId)

	return nil
}

type ProductCreationResponse struct {
	Result int `json:"result"`
}

func AddProductsWithImage(product Product, sectionID string) (int, error) {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "crm.product.add"

	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	requestBody := map[string]interface{}{
		"fields": map[string]interface{}{
			"NAME":       product.Name,
			"PRICE":      product.Price,
			"SECTION_ID": sectionID,
			"PROPERTY_197": []map[string]interface{}{
				{
					"valueId":  0,
					"fileData": []string{"avatar.png", product.ImageBase64},
				},
			},
			"PROPERTY_198": product.Material,
			"PROPERTY_199": product.Laser,
			"PROPERTY_200": product.Bend,
			"PROPERTY_201": product.Paint,
			"PROPERTY_202": product.Production,
			"PROPERTY_203": product.PipeCutting,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return 0, fmt.Errorf("error marshalling request body: %v", err)
	}

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

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	var response ProductCreationResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return 0, fmt.Errorf("error unmarshalling response: %v", err)
	}

	if response.Result == 0 {
		return 0, fmt.Errorf("failed to create product, response: %s", string(responseData))
	}

	log.Println("Product added with ID:", response.Result)
	return response.Result, nil
}

func processLaser(orderNumber, client, fileName string, smartProcessID, engineerID int) (int, error) {
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
	taskID, err := AddTaskToGroup(orderNumber, client, "Производство", engineerID, 2, 1046, smartProcessID)
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

/*func HandlerAddCustomTaskToParentId(w http.ResponseWriter, r *http.Request) {
	products, err := processProducts("file.xlsx", 688, 149)
	if err != nil {
		log.Printf("Error processing products: %v\n", err)
	}
	fmt.Fprintf(w, "Products processed successfully: %v", products)
}*/

func ConductDocumentId(documentID int) error {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "catalog.document.conduct"

	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	// Формируем тело запроса
	requestBody := map[string]interface{}{
		"id": documentID, // ID документа, передается в теле запроса
	}

	// Преобразуем тело запроса в JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %v", err)
	}

	// Создаем HTTP запрос
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	// Парсим ответ
	var response map[string]interface{}
	if err := json.Unmarshal(responseData, &response); err != nil {
		return fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Проверяем ошибки в ответе
	if _, ok := response["error"]; ok {
		log.Printf("Full response on error: %v", string(responseData))
		return fmt.Errorf("ошибка: %s", response["error_description"])
	}

	// Проверяем наличие результата
	if _, ok := response["result"]; !ok {
		log.Printf("Unexpected response structure: %v", string(responseData))
		return fmt.Errorf("missing 'result' in response")
	}

	log.Printf("Документ успешно проведен: ID документа: %d\n", documentID)

	return nil
}

func ReadXlsProducts(filename string) map[string][]string {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		log.Println("Error opening file:", err)
		return nil
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("Error closing file:", err)
		}
	}()

	// Для хранения данных, где ключ — значение из первой ячейки строки
	data := make(map[string][]string)

	// Получаем все строки из листа "КП"
	rows, err := f.GetRows("КП")
	if err != nil {
		log.Println("Error getting rows:", err)
		return nil
	}

	// Перебираем строки
	for rowIndex, row := range rows {
		// Пропускаем первую строку (например, заголовки)
		if rowIndex == 0 {
			continue
		}

		// Если строка пустая, пропускаем её
		if len(row) == 0 {
			continue
		}

		// Первый элемент строки становится ключом
		key := row[0]
		// Остальные элементы добавляются в значение
		if len(row) > 1 {
			data[key] = row[1:]
		} else {
			data[key] = []string{}
		}
	}

	return data
}

// Response структура для обработки ответа от Bitrix24
type DealResponse struct {
	Result struct {
		UFCRM1733146336 int `json:"UF_CRM_1733146336"`
	} `json:"result"`
}

// GetDealFieldValue функция для получения значения UF_CRM_1733146336 по ID сделки
func GetProductionEngineerIdByDeal(dealID string) (int, error) {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "crm.deal.get"
	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	// Формируем тело запроса
	requestBody := map[string]interface{}{
		"id": dealID,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return 0, fmt.Errorf("error marshalling request body: %v", err)
	}

	// Создаем HTTP-запрос
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
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

	// Разбираем JSON-ответ
	var response DealResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return 0, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Проверяем, есть ли поле в ответе
	if response.Result.UFCRM1733146336 == 0 {
		return 0, fmt.Errorf("field UF_CRM_1733146336 not found or is empty in the response")
	}

	log.Printf("Field UF_CRM_1733146336 value: %s\n", response.Result.UFCRM1733146336)
	return response.Result.UFCRM1733146336, nil
}

/*func GetInfoAboutProductsFields() {
	//https://bitrix.laser-flex.ru/rest/149/ycz7102vaerygxvb/profile.json

	bitrixMethod := "crm.item.productrow.fields"

	requestURL := fmt.Sprintf("%s/rest/%s?", endpoints.BitrixDomain, bitrixMethod)

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

}*/
