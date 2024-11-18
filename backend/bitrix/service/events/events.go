package events

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func OnCrmDealAddEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("OnCrmDealAddEvent started: ")
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading request body:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	log.Println("resp_at_first:", string(bs))
	defer r.Body.Close()

	values, err := url.ParseQuery(string(bs))
	if err != nil {
		panic(err)
	}

	// Convert to a map (manually handling nested structures)
	result := make(map[string]interface{})
	for key, value := range values {
		// Assuming only one value per key for simplicity
		result[key] = value[0]
	}

	// Converting the map to JSON
	jsonBytes, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonBytes))

	// To extract the ID, you can unmarshal the relevant part of the JSON back into a Go data structure
	// or manipulate the JSON directly if you prefer
	type DataFields struct {
		ID string `json:"data[FIELDS][ID]"`
	}

	var dataFields DataFields
	if err := json.Unmarshal(jsonBytes, &dataFields); err != nil {
		panic(err)
	}

	fmt.Println("Extracted ID:", dataFields.ID)

}

//func OnCrmDealUpdateEvent(w http.ResponseWriter, r *http.Request) {}
//func OnCrmDealDeleteEvent(w http.ResponseWriter, r *http.Request) {}
