package companies

import (
	"bitrix_app/backend/bitrix/authorize"
	"bitrix_app/backend/bitrix/endpoints"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// GetCompany retrieves the company data by ID
func GetCompany(companyID string) (*Company, error) {
	bitrixMethod := "crm.company.get"
	requestURL := fmt.Sprintf("%s/rest/%s?auth=%s&id=%s", endpoints.BitrixDomain, bitrixMethod, authorize.GlobalAuthId, companyID)

	// Create a new request
	req, err := http.NewRequest("GET", requestURL, nil)
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
	//log.Println("GetCompany Response:", string(responseData))

	// Parse the response into a Company struct
	var response struct {
		Result Company `json:"result"`
	}
	if err := json.Unmarshal(responseData, &response); err != nil {
		log.Println("Error unmarshaling response:", err)
		return nil, err
	}

	return &response.Result, nil
}

// CompanyHandler handles the request for company data
func CompanyHandler(w http.ResponseWriter, r *http.Request) {
	companyID := r.URL.Query().Get("id")
	if companyID == "" {
		http.Error(w, "Company ID is missing", http.StatusBadRequest)
		return
	}

	company, err := GetCompany(companyID)
	if err != nil {
		log.Println("Error getting company:", err)
		http.Error(w, "Error getting company", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type to application/json and write the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}
