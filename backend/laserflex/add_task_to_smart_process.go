package laserflex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func AddCustomTaskToParentId(title string, responsibleID, groupID int, customFields CustomTaskFields, elementID int) (int, error) {
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
			"UF_AUTO_303168834495": []string{customFields.OrderNumber}, // № заказа
			"UF_AUTO_876283676967": []string{customFields.Customer},    // Заказчик
			"UF_AUTO_794809224848": []string{customFields.Manager},     // Менеджер
			"UF_AUTO_468857876599": []string{customFields.Material},    // Материал
			"UF_AUTO_497907774817": []string{customFields.Comment},     // Комментарий
			"UF_AUTO_552243496167": []string{customFields.Quantity},    // Кол-во
			"DEADLINE":             deadline,                           // DEADLINE: текущая дата + 13 часов
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
