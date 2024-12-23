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

func AddCustomTaskToParentId(title string, responsibleID, groupID int, customFields CustomTaskFields, elementID int, timeEstimate int) (int, error) {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "tasks.task.add"
	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	// Генерация ссылки смарт-процесса
	smartProcessLink := GenerateSmartProcessLink(1046, elementID)

	// Подготовка тела запроса
	requestBody := map[string]interface{}{
		"fields": map[string]interface{}{
			"TITLE":                title,
			"RESPONSIBLE_ID":       responsibleID,
			"GROUP_ID":             groupID,
			"UF_CRM_TASK":          []string{smartProcessLink},
			"UF_AUTO_303168834495": []string{customFields.OrderNumber}, // № заказа
			"UF_AUTO_876283676967": []string{customFields.Customer},    // Заказчик
			"UF_AUTO_794809224848": []string{customFields.Manager},     // Менеджер
			"UF_AUTO_468857876599": []string{customFields.Material},    // Материал
			"UF_AUTO_497907774817": []string{customFields.Comment},     // Комментарий
			"UF_AUTO_552243496167": []string{customFields.Quantity},    // Кол-во
			"DEADLINE":             time.Now().Add(16 * time.Hour).Format("02.01.2006T15:04:05"),
			"ALLOW_TIME_TRACKING":  "Y",
			"TIME_ESTIMATE":        timeEstimate,
		},
	}

	// Сериализация тела запроса
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

	// Читаем ответ
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	log.Printf("Attention!!!!!!!!! AddCustomTaskToParentId Response from Bitrix24: %s\n", string(responseData))

	// Разбираем ответ
	var response TaskResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return 0, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Обрабатываем ID задачи
	var taskID int
	if err := json.Unmarshal(response.Result.Task.ID, &taskID); err != nil {
		// Если ID строка, конвертируем в число
		var taskIDStr string
		if err := json.Unmarshal(response.Result.Task.ID, &taskIDStr); err != nil {
			return 0, fmt.Errorf("error parsing task id: %v", err)
		}
		taskID, err = strconv.Atoi(taskIDStr)
		if err != nil {
			return 0, fmt.Errorf("error converting task id to int: %v", err)
		}
	}

	if taskID == 0 {
		return 0, fmt.Errorf("failed to create subtask, response: %s", string(responseData))
	}

	log.Printf("Subtask created with ID: %d\n", taskID)
	return taskID, nil
}

func processLaser(fileName string, smartProcessID, engineerID int) (int, error) {
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

func HandlerAddCustomTaskToParentId(w http.ResponseWriter, r *http.Request) {
	products, err := processProducts("file.xlsx", 688, 149)
	if err != nil {
		log.Printf("Error processing products: %v\n", err)
	}
	fmt.Fprintf(w, "Products processed successfully: %v", products)
}
