package laserflex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func AddCatalogDocument(dealID string, assignedById int, totalProductsPrice float64) (int, error) {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "catalog.document.add"

	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	// Формируем тело запроса
	requestBody := map[string]interface{}{
		"fields": map[string]interface{}{
			"docType":       "A",
			"responsibleId": assignedById,
			"createdBy":     assignedById,
			"currency":      "RUB",
			"status":        "S",
			"statusBy":      assignedById,
			"total":         fmt.Sprintf("%.2f", totalProductsPrice), // Округляем до двух знаков
			"title":         dealID,
			"contractorId":  556,
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
