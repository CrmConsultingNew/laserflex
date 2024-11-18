package requisites

import (
	"bitrix_app/backend/bitrix/endpoints"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetRequisitesByCompanyID(companyId string, authKey string) ([]Requisites, error) {
	bitrixMethod := "crm.requisite.list"
	var requisites []Requisites

	for {
		requestURL := fmt.Sprintf("%s/rest/%s?auth=%s", endpoints.BitrixDomain, bitrixMethod, authKey)

		// Формируем правильное тело запроса
		requestBody := fmt.Sprintf(`{
            "filter": {
                "ENTITY_ID": "%s"
            }
        }`, companyId)

		// Здесь не нужно использовать json.Marshal, так как requestBody уже является корректной JSON-строкой
		req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer([]byte(requestBody)))
		if err != nil {
			log.Println("Error creating new request:", err)
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")

		// Отправляем запрос
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Error sending request:", err)
			return nil, err
		}
		defer resp.Body.Close()

		// Читаем ответ
		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
			return nil, err
		}

		// Логирование для отладки
		log.Println("GetRequisitesByCompanyID Response:", string(responseData))

		// Парсим ответ
		var response struct {
			Result []Requisites `json:"result"`
			Next   int          `json:"next"`
		}
		if err := json.Unmarshal(responseData, &response); err != nil {
			log.Println("Error unmarshaling response:", err)
			return nil, err
		}

		requisites = append(requisites, response.Result...)

		// Если нет следующей страницы (next), завершаем цикл
		if response.Next == 0 {
			break
		}

		// Увеличиваем offset для получения следующей порции данных (если требуется)
		// start = response.Next
	}

	return requisites, nil
}
