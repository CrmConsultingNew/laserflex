package description

import (
	"bitrix_app/backend/bitrix/authorize"
	"bitrix_app/backend/bitrix/endpoints"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func DescriptionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	pathSegments := strings.Split(r.URL.Path, "/")

	// Assuming the URL format is /api/documents/{ID}
	// The ID is expected to be the fourth segment, hence index 3 (0-based index)
	if len(pathSegments) < 4 {
		http.Error(w, "Invalid request path", http.StatusBadRequest)
		return
	}
	entityId := pathSegments[3]

	// Use the extracted ID as needed, for now, we'll just print it
	fmt.Println("Extracted ID DescriptionHandler: ", entityId)

	docs, err := GetDescription(authorize.GlobalAuthId, entityId)
	if err != nil {
		log.Println("Error getting documentation in DescriptionHandler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Directly write the mock JSON as the response
	_, err = w.Write(docs) // docs already contains the JSON
	if err != nil {
		log.Println("Error writing the mock response:", err)
	}
}

func GetDescription(authID string, dealId string) ([]byte, error) {
	bitrixMethod := "crm.deal.get"

	requestURL := fmt.Sprintf("%s/rest/%s?auth=%s", endpoints.BitrixDomain, bitrixMethod, authID)

	// The body here needs to be an object that matches the expected structure for the Bitrix24 API call
	bodyObj := map[string]string{"id": dealId}
	jsonData, err := json.Marshal(bodyObj)
	if err != nil {
		log.Println("Error marshaling request body:", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating new request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json") // Ensure the content type is set to application/json

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

	// Assuming the response will be a single deal
	var dealResponse DealResponse
	if err := json.Unmarshal(responseData, &dealResponse); err != nil {
		log.Printf("Error unmarshalling response into DealResponse: %v", err)
		return nil, err
	}

	// If needed, marshal the DealResponse (or just the Deal part of it) back into JSON to return
	jsonResponse, err := json.Marshal(dealResponse)
	if err != nil {
		log.Printf("Error marshalling DealResponse back to JSON: %v", err)
		return nil, err
	}
	fmt.Println("string(jsonResponse): ", string(jsonResponse))

	return jsonResponse, nil
}
