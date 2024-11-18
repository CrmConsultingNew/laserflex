package products

import (
	"bitrix_app/backend/bitrix/endpoints"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func AddProduct(name string, currency string, price float64, sort int, authID string) error {
	bitrixMethod := "crm.product.add"
	requestURL := fmt.Sprintf("%s/rest/%s?auth=%s", endpoints.BitrixDomain, bitrixMethod, authID)

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
}

func AddMultipleProducts(authID string) {
	products := []struct {
		Name     string
		Currency string
		Price    float64
		Sort     int
	}{
		{"1С-Битрикс: Управление сайтом - Старт", "RUB", 4900, 500},
		{"1С-Битрикс: Управление сайтом - Премиум", "RUB", 14900, 600},
		{"1С-Битрикс: CRM для бизнеса", "RUB", 8900, 700},
	}

	for _, product := range products {
		if err := AddProduct(product.Name, product.Currency, product.Price, product.Sort, authID); err != nil {
			log.Println("Ошибка при добавлении товара:", err)
		}
	}
}
