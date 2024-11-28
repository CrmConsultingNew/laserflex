package laserflex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

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
