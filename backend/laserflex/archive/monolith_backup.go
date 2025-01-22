package archive

import (
	"bitrix_app/backend/laserflex"
	"bitrix_app/backend/laserflex/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"net/http"
	"strings"
)

/*func HandlerProcessProducts(w http.ResponseWriter, r *http.Request) {
	products, err := processProducts("file.xlsx", 688, 149)
	if err != nil {
		log.Printf("Error processing products: %v\n", err)
	}
	fmt.Fprintf(w, "Products processed successfully: %v", products)
}*/

// processTask универсальная функция для обработки задач
/*func processTask(orderNumber, client, fileName string, smartProcessID, engineerID int, taskType string, groupID int) (int, error) {
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
}*/

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

			customFields := models.CustomTaskFields{
				OrderNumber: orderNumber,
				Customer:    customer,
			}

			// Создаём задачу, передавая массив уникальных цветов
			colorArray := make([]string, 0, len(uniqueColors))
			for color := range uniqueColors {
				colorArray = append(colorArray, color)
			}

			_, err := laserflex.AddCustomCoatingTask(taskTitle, engineerID, 12, customFields, smartProcessID, colorArray)
			if err != nil {
				log.Printf("Error creating coating task: %v\n", err)
				continue
			}
		}
	}

	log.Printf("Processing coating tasks completed successfully")
	return nil
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
		if value := laserflex.getValue(cells, headers["Лазерные работы"]); value != "" {
			parsedRow.LaserWorks = &WorkGroup{
				GroupID: 1,
				Data: laserflex.extractData(cells, headers, []string{
					"№ заказа", "Заказчик", "Менеджер", "Количество материала", "Лазерные работы", "Нанесение покрытий", "Комментарий",
				}),
			}
		}

		if value := laserflex.getValue(cells, headers["Труборез"]); value != "" {
			parsedRow.PipeCutting = &WorkGroup{
				GroupID: 11,
				Data: laserflex.extractData(cells, headers, []string{
					"№ заказа", "Заказчик", "Менеджер", "Количество материала", "Труборез", "Нанесение покрытий", "Комментарий",
				}),
			}
		}

		if value := laserflex.getValue(cells, headers["Гибочные работы"]); value != "" {
			parsedRow.BendWorks = &WorkGroup{
				GroupID: 10,
				Data: laserflex.extractData(cells, headers, []string{
					"№ заказа", "Заказчик", "Менеджер", "Количество материала", "Гибочные работы", "Нанесение покрытий", "Комментарий",
				}),
			}
		}

		if value := laserflex.getValue(cells, headers["Производство"]); value != "" {
			parsedRow.Production = &WorkGroup{
				GroupID: 2,
				Data: laserflex.extractData(cells, headers, []string{
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

/*func processLaser(orderNumber, client, fileName string, smartProcessID, engineerID int) (int, error) {
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
}*/

/*func HandlerAddCustomTaskToParentId(w http.ResponseWriter, r *http.Request) {
	products, err := processProducts("file.xlsx", 688, 149)
	if err != nil {
		log.Printf("Error processing products: %v\n", err)
	}
	fmt.Fprintf(w, "Products processed successfully: %v", products)
}*/

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
