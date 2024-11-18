package settings

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type InputFields struct {
	InputField1 string `json:"input_field1"`
	InputField2 string `json:"input_field2"`
}

func SaveSettingsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	var inputs []InputFields
	err := json.NewDecoder(r.Body).Decode(&inputs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert the slice to a map
	result := make(map[string]InputFields)
	for i, input := range inputs {
		// Using i+1 to start numbering from 1 instead of 0
		key := fmt.Sprintf("%d", i+1)
		result[key] = input
	}

	// For demonstration, printing the map to the server console
	fmt.Printf("%+v\n", result)

	// Respond back to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}
