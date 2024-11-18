package smart_processes

type ItemListRequest struct {
	CompanyID string `json:"companyID"`
}

/*func ItemListHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ItemListHandler was started and ready")

	body, _ := io.ReadAll(r.Body)
	log.Println("body ItemsListHandler: ", body)

	companyID := r.URL.Query().Get("id")
	if companyID == "" {
		http.Error(w, "Company ID is missing", http.StatusBadRequest)
		return
	}

	// Вызов функции GetItemByCompany с переданным companyID
	itemsData, err := GetItemByCompany(companyID)
	if err != nil {
		http.Error(w, "Failed to get processes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the Content-Type to application/json and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(itemsData)
	if err != nil {
		log.Println("Failed to encode items data to json: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("ItemListHandler was ended")
}*/
