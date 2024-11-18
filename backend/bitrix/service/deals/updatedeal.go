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

func UpdateDeal(dealID string, deal DealInfo, authID string) (*DealInfo, error) {
	bitrixMethod := "crm.deal.update"
	requestURL := fmt.Sprintf("%s/rest/%s?auth=%s", endpoints.BitrixDomain, bitrixMethod, authID)

	// Создаем тело запроса с полями в нужном формате
	requestBody := map[string]interface{}{
		"id": dealID,
		"fields": map[string]interface{}{
			"TITLE":             deal.Title,
			"UF_CRM_1726587543": deal.NISHA,
			"UF_CRM_1726587678": deal.GEOGRAPHIA,
			"UF_CRM_1726587712": deal.KONECHNIY_PRODUCT,
			"UF_CRM_1706503695": deal.KOLICHESTVO_SOTRUDNIKOV,
			"UF_CRM_1726587982": deal.KOLICHESTVO_SOTRUDNIKOV_OP,
			"UF_CRM_1726588010": deal.EST_ROP,
			"UF_CRM_1726588094": deal.EST_HR,
			"UF_CRM_1726588113": deal.OSNOVNIE_KANALI_PRODAJ,
			"UF_CRM_1726588197": deal.SEZONNOST,
			"UF_CRM_1726588214": deal.ZAPROS,
			"UF_CRM_1726588301": deal.KAKIE_CELI_PERED_KOMPANIEI,
			"UF_CRM_1726588336": deal.CHTO_OJIDAETE_OT_SOTRUDNICHESTVA,
			"COMPANY_ID":        deal.CompanyID, // Добавляем CompanyID
		},
	}

	// Преобразуем тело запроса в JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Println("Error marshaling request body:", err)
		return nil, err
	}

	// Выводим тело запроса в консоль
	log.Println("Request JSON data:", string(jsonData))

	// Создаем новый POST запрос с телом
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
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

	// Чтение и обработка тела ответа
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil, err
	}

	// Выводим полный ответ в консоль
	log.Println("Response data:", string(responseData))

	// Разбор ответа
	var response struct {
		Result bool   `json:"result"`
		Error  string `json:"error,omitempty"`
		Time   struct {
			Start            float64 `json:"start"`
			Finish           float64 `json:"finish"`
			Duration         float64 `json:"duration"`
			Processing       float64 `json:"processing"`
			DateStart        string  `json:"date_start"`
			DateFinish       string  `json:"date_finish"`
			OperatingResetAt float64 `json:"operating_reset_at"`
			Operating        float64 `json:"operating"`
		} `json:"time"`
	}

	if err := json.Unmarshal(responseData, &response); err != nil {
		log.Println("Error unmarshaling response:", err)
		return nil, err
	}

	if response.Error != "" {
		log.Println("Bitrix24 API error:", response.Error)
		return nil, fmt.Errorf("API error: %s", response.Error)
	}

	log.Println("Response result:", response.Result)
	log.Printf("Time info: %+v\n", response.Time)

	// Возвращаем nil, так как в ответе нет данных о сделке
	return nil, nil
}
