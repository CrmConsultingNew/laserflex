package deals

import (
	"bitrix_app/backend/bitrix/endpoints"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetInfoAboutDealByID(dealID string, authID string) (*DealInfo, error) {
	// Формируем URL для API запроса
	requestUrl := fmt.Sprintf("%s/rest/crm.deal.get?auth=%s", endpoints.BitrixDomain, authID)

	// Подготавливаем тело запроса
	jsonData := fmt.Sprintf(`{"id": "%s"}`, dealID)

	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		log.Println("error creating new request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response body:", err)
		return nil, err
	}

	// Парсим ответ в структуру
	var apiResponse ApiResponse
	err = json.Unmarshal(bz, &apiResponse)
	if err != nil {
		log.Println("error unmarshalling response to ApiResponse:", err)
		return nil, err
	}

	dealInfo := apiResponse.Result
	log.Println("dealInfo 1st request:", dealInfo)
	return &dealInfo, nil
}
