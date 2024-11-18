package smart_processes

import (
	"bitrix_app/backend/bitrix/authorize"
	"bitrix_app/backend/bitrix/endpoints"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Filter struct {
	CompanyID string `json:"companyId"`
}

type ItemsRequestBody struct {
	EntityTypeId int    `json:"entityTypeId"`
	Filter       Filter `json:"filter"`
}

// ItemsResponse Акты поступлений (смарт процесс)
type ItemsResponse struct {
	ID              int     `json:"id"`
	Opportunity     float64 `json:"opportunity"`
	UFCrm1712128088 string  `json:"ufCrm26_1712128088"` // Дата Акта
}

func formatDate(dateStr string) (string, error) {
	parsedDate, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return "", err
	}
	return parsedDate.Format("02.01.2006"), nil
}

func GetItemByCompany(companyID string) ([]ItemsResponse, error) {
	bitrixMethod := "crm.item.list"
	requestURL := fmt.Sprintf("%s/rest/%s?auth=%s", endpoints.BitrixDomain, bitrixMethod, authorize.GlobalAuthId)

	body := ItemsRequestBody{
		EntityTypeId: 179,
		Filter: Filter{
			CompanyID: companyID,
		},
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		log.Println("Error marshaling request body:", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating new request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil, err
	}

	log.Println("GetItemByCompany Response:", string(responseData))

	var rawResponse struct {
		Result struct {
			Items []struct {
				ID              int     `json:"ID"`
				Opportunity     float64 `json:"OPPORTUNITY"`
				UFCrm1712128088 string  `json:"ufCrm26_1712128088"`
			} `json:"items"`
		} `json:"result"`
	}
	if err := json.Unmarshal(responseData, &rawResponse); err != nil {
		log.Println("Error unmarshaling response:", err)
		return nil, err
	}

	log.Printf("Total items received: %d", len(rawResponse.Result.Items))

	var items []ItemsResponse
	for _, item := range rawResponse.Result.Items {
		// Форматируем дату
		formattedDate, err := formatDate(item.UFCrm1712128088)
		if err != nil {
			log.Printf("Error formatting date: %v", err)
			formattedDate = item.UFCrm1712128088 // Если произошла ошибка, оставляем оригинальную дату
		}

		log.Printf("Item ID: %d, Opportunity: %f, UF_CRM_26_171212808: %s", item.ID, item.Opportunity, formattedDate)
		items = append(items, ItemsResponse{
			ID:              item.ID,
			Opportunity:     item.Opportunity,
			UFCrm1712128088: formattedDate,
		})
	}

	if len(items) == 0 {
		log.Printf("No opportunities found for company ID: %s", companyID)
	}

	return items, nil
}

// GetItemsByCompanyHandler передает значения item.Opportunity и item.UFCrm1712128088 на фронтенд
func GetItemsByCompanyHandler(w http.ResponseWriter, r *http.Request) {
	companyID := r.URL.Query().Get("id")
	if companyID == "" {
		http.Error(w, "Company ID is missing", http.StatusBadRequest)
		return
	}

	// Получаем данные из GetItemByCompany
	items, err := GetItemByCompany(companyID)
	if err != nil {
		http.Error(w, "Error fetching items data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Преобразуем данные в JSON для отправки на фронтенд
	responseData, err := json.Marshal(items)
	if err != nil {
		http.Error(w, "Error marshaling items data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки и отправляем данные на фронтенд
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}
