package laserflex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

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
