package laserflex

/*func GetInfoAboutProductsFields() {
	//https://bitrix.laser-flex.ru/rest/149/ycz7102vaerygxvb/profile.json

	bitrixMethod := "crm.item.productrow.fields"

	requestURL := fmt.Sprintf("%s/rest/%s?", endpoints.BitrixDomain, bitrixMethod)

	requestBody := map[string]interface{}{
		"fields": map[string]interface{}{
			"NAME":        name,
			"CURRENCY_ID": currency,
			"PRICE":       price,
			"SORT":        sort,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(responseData, &response); err != nil {
		return err
	}

	if _, ok := response["error"]; ok {
		return fmt.Errorf("Ошибка: %s", response["error_description"])
	}

	log.Println("Товар добавлен с ID:", response["result"])
	return nil

}*/
