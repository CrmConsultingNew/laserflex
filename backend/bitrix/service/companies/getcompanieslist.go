package companies

import (
	"bitrix_app/backend/bitrix/authorize"
	"bitrix_app/backend/bitrix/endpoints"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func CompaniesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Service CompaniesHandler was started")

	// Получение списка компаний
	companiesData, err := GetAllCompaniesList(authorize.GlobalAuthId)
	if err != nil {
		log.Println("Error getting companies:", err)
		http.Error(w, "Error getting companies", http.StatusInternalServerError)
		return
	}

	// Преобразование данных в JSON
	jsonData, err := json.Marshal(companiesData)
	if err != nil {
		log.Println("Error marshaling companies data:", err)
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	// Установка заголовков и отправка данных клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func GetCompanyByID(companyID, authKey string) (Company, error) {
	bitrixMethod := "crm.company.get"

	// Формируем URL запроса к API Bitrix
	requestURL := fmt.Sprintf("%s/rest/%s?auth=%s", endpoints.BitrixDomain, bitrixMethod, authKey)

	// Формируем тело запроса
	reqBody := fmt.Sprintf(`{"id": "%s"}`, companyID)
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		log.Println("Error creating new request:", err)
		return Company{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос к API
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return Company{}, err
	}
	defer resp.Body.Close()

	// Читаем и обрабатываем ответ
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return Company{}, err
	}

	// Логируем для отладки
	log.Println("GetCompanyByID Response:", string(responseData))

	// Разбираем ответ в структуру Company
	var response struct {
		Result Company `json:"result"`
	}
	if err := json.Unmarshal(responseData, &response); err != nil {
		log.Println("Error unmarshaling response:", err)
		return Company{}, err
	}

	return response.Result, nil
}

// GetAllCompaniesList retrieves all companies from Bitrix by calling crm.company.list multiple times until all data is retrieved.
func GetAllCompaniesList(authKey string) ([]Company, error) {
	bitrixMethod := "crm.company.list"
	allCompanies := []Company{}
	start := 0

	for {
		requestURL := fmt.Sprintf("%s/rest/%s?auth=%s", endpoints.BitrixDomain, bitrixMethod, authKey)

		// Prepare request body with pagination
		requestBody := map[string]interface{}{
			"start": start, // Начало выборки
		}

		// Marshal the request body into JSON
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			log.Println("Error marshaling request body:", err)
			return nil, err
		}

		// Create a new request with JSON body
		req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("Error creating new request:", err)
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")

		// Send the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Error sending request:", err)
			return nil, err
		}
		defer resp.Body.Close()

		// Read and parse the response body
		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
			return nil, err
		}

		// Log the raw response for debugging
		//log.Println("GetCompaniesList Response:", string(responseData))

		// Parse the response into a slice of companies
		var response struct {
			Result []Company `json:"result"`
			Next   int       `json:"next"`
		}
		if err := json.Unmarshal(responseData, &response); err != nil {
			log.Println("Error unmarshaling response:", err)
			return nil, err
		}

		allCompanies = append(allCompanies, response.Result...)

		// Если поле next пустое или равно 0, завершить цикл
		if response.Next == 0 {
			break
		}

		start = response.Next
	}

	return allCompanies, nil
}
