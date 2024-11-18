package companies

import (
	"bitrix_app/backend/bitrix/endpoints"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func UpdateCompany(companyID string, company Company, authID string) (*Company, error) {
	bitrixMethod := "crm.company.update"
	requestURL := fmt.Sprintf("%s/rest/%s?auth=%s", endpoints.BitrixDomain, bitrixMethod, authID)

	// Создаем тело запроса с полями в нужном формате
	requestBody := map[string]interface{}{
		"id": companyID,
		"fields": map[string]interface{}{
			"TITLE":             company.Title,
			"UF_CRM_1726674632": company.INN,
			"UF_CRM_1726587588": company.NISHA,
			"UF_CRM_1726587646": company.GEOGRAPHIA,
			"UF_CRM_1726587741": company.KONECHNIY_PRODUCT,
			"UF_CRM_1726587790": company.OBOROT_KOMPANII,
			"UF_CRM_1726587900": company.SREDNEMESYACHNAYA_VIRYCHKA,
			"UF_CRM_EMPLOYEES":  company.KOLICHESTVO_SOTRUDNIKOV,
			"UF_CRM_1726587932": company.KOLICHESTVO_SOTRUDNIKOV_OP,
			"UF_CRM_1726588037": company.EST_ROP,
			"UF_CRM_1726588054": company.EST_HR,
			"UF_CRM_1726588155": company.OSNOVNIE_KANALI_PRODAJ,
			"UF_CRM_1726588171": company.SEZONNOST,
			"UF_CRM_1726588242": company.ZAPROS,
			"UF_CRM_1726588284": company.KAKIE_CELI_PERED_KOMPANIEI,
			"UF_CRM_1726588358": company.CHTO_OJIDAETE_OT_SOTRUDNICHESTVA,
			"ASSIGNED_BY_ID":    company.AssignedByID,
			"HAS_EMAIL":         company.HasEmail,
			"EMAIL":             company.Emails,
		},
	}

	// Преобразуем тело запроса в JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Println("Error marshaling request body:", err)
		return nil, err
	}

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

	// Разбор ответа в структуру Company
	var response struct {
		Result json.RawMessage `json:"result"`
		Error  string          `json:"error,omitempty"`
	}
	if err := json.Unmarshal(responseData, &response); err != nil {
		log.Println("Error unmarshaling response:", err)
		return nil, err
	}

	if response.Error != "" {
		log.Println("Bitrix24 API error:", response.Error)
		return nil, fmt.Errorf("API error: %s", response.Error)
	}

	var companyResult Company
	if err := json.Unmarshal(response.Result, &companyResult); err != nil {
		log.Println("Error unmarshaling company result:", err)
		return nil, err
	}

	log.Println("response.Result: ", companyResult)
	return &companyResult, nil
}
