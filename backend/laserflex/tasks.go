package laserflex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// AddTaskToGroup создает задачу и возвращает ID созданной задачи
func AddTaskToGroup(title string, responsibleID, groupID, processTypeID, elementID int) (int, error) {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "tasks.task.add"
	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	// Генерация ссылки смарт-процесса
	smartProcessLink := GenerateSmartProcessLink(processTypeID, elementID)

	// Формирование тела запроса
	requestBody := TaskRequest{
		Fields: map[string]interface{}{
			"TITLE":          title,
			"RESPONSIBLE_ID": responsibleID,
			"GROUP_ID":       groupID,
			"UF_CRM_TASK":    []string{smartProcessLink},
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

	// Логирование для отладки
	log.Printf("Raw Response: %s\n", string(responseData))

	// Измененная структура ответа
	var response struct {
		Result struct {
			Task struct {
				ID json.RawMessage `json:"id"` // Используем RawMessage для обработки id как строки или числа
			} `json:"task"`
		} `json:"result"`
	}

	if err := json.Unmarshal(responseData, &response); err != nil {
		return 0, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Обрабатываем ID задачи
	var taskID int
	if err := json.Unmarshal(response.Result.Task.ID, &taskID); err != nil {
		// Если id не является числом, пробуем распарсить его как строку
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
		return 0, fmt.Errorf("failed to create task, response: %s", string(responseData))
	}

	log.Printf("Task created with ID: %d\n", taskID)
	return taskID, nil
}

// GenerateSmartProcessLink генерирует идентификатор смарт-процесса
func GenerateSmartProcessLink(processTypeID, elementID int) string {
	// Преобразуем идентификатор типа смарт-процесса в шестнадцатеричную систему
	hexType := strings.ToLower(strconv.FormatInt(int64(processTypeID), 16))

	// Формируем строку привязки
	return fmt.Sprintf("T%s_%d", hexType, elementID)
}

// ID группы 1 - Лазерные работы
// ID группы 11 - Труборез
// ID группы 10 - Гибочные работы
// ID группы 2 - Производство

// AddTaskToParentId создает подзадачу с привязкой к PARENT_ID и пользовательскими полями
func AddTaskToParentId(title string, responsibleID, groupID int, customFields CustomTaskFields) (int, error) {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "tasks.task.add"
	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	// Подготовка тела запроса
	requestBody := map[string]interface{}{
		"fields": map[string]interface{}{
			"TITLE":          title,
			"RESPONSIBLE_ID": responsibleID,
			"GROUP_ID":       groupID,
			// Добавляем пользовательские поля с учетом формата массива
			"UF_AUTO_303168834495": []string{customFields.OrderNumber},    // № заказа
			"UF_AUTO_876283676967": []string{customFields.Customer},       // Заказчик
			"UF_AUTO_794809224848": []string{customFields.Manager},        // Менеджер
			"UF_AUTO_468857876599": []string{customFields.Material},       // Материал
			"UF_AUTO_497907774817": []string{customFields.Comment},        // Комментарий
			"UF_AUTO_433735177517": []string{customFields.ProductionTask}, // Произв. Задача
			"UF_AUTO_726724682983": []string{customFields.Bend},           // Гибка
			"UF_AUTO_512869473370": []string{customFields.Coating},        // Покрытие
			"UF_AUTO_555642596740": customFields.TemporaryOrderSum,        // Временная сумма заказа (одиночное значение)
			"UF_AUTO_552243496167": []string{customFields.Quantity},       // Кол-во
		},
	}

	// Сериализация тела запроса
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

	// Логируем ответ для отладки
	log.Printf("Response from Bitrix24: %s\n", string(responseData))

	// Разбираем ответ
	var response TaskResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return 0, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Проверяем успешность создания подзадачи
	if response.Result.Task.ID == 0 {
		return 0, fmt.Errorf("failed to create subtask, response: %s", string(responseData))
	}

	log.Printf("Subtask created with ID: %d\n", response.Result.Task.ID)
	return response.Result.Task.ID, nil
}

// AddTaskWithChecklist создает задачу с чек-листом и возвращает ID созданной задачи
func AddCheckListToTheTask(taskID int, title string) (int, error) {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "task.checklistitem.add"
	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	// Формирование тела запроса
	requestBody := []interface{}{
		taskID,
		map[string]interface{}{
			"TITLE":       title,
			"IS_COMPLETE": "N",
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

	var response struct {
		Result struct {
			ID int `json:"id"`
		} `json:"result"`
	}
	if err := json.Unmarshal(responseData, &response); err != nil {
		return 0, fmt.Errorf("error unmarshalling response: %v", err)
	}

	if response.Result.ID == 0 {
		return 0, fmt.Errorf("failed to add checklist item, response: %s", string(responseData))
	}

	log.Printf("Checklist item added with ID: %d\n", response.Result.ID)
	return response.Result.ID, nil
}
