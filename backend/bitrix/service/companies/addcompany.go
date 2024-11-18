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

func AddCompany(company Company, authID string) (int, error) {
	bitrixMethod := "crm.company.add"
	requestURL := fmt.Sprintf("%s/rest/%s?auth=%s", endpoints.BitrixDomain, bitrixMethod, authID)

	// Создаем тело запроса с полями из структуры Company
	requestBody := map[string]interface{}{
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
		return 0, err
	}

	// Создаем новый POST запрос с телом
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating new request:", err)
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return 0, err
	}
	defer resp.Body.Close()

	// Чтение и обработка тела ответа
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return 0, err
	}

	log.Println("response data:", string(responseData))
	// Логирование сырых данных ответа для отладки
	//log.Println("AddCompany Response:", string(responseData))

	// Разбор ответа в структуру Company
	var response struct {
		Result int `json:"result"`
	}

	// Разбираем ответ в новую структуру
	if err := json.Unmarshal(responseData, &response); err != nil {
		log.Println("Error unmarshaling response:", err)
		return 0, err
	}

	// Теперь вы можете получить значение result
	resultValue := response.Result
	//log.Println("Extracted result value:", resultValue) // Это будет 950

	return resultValue, nil
}
