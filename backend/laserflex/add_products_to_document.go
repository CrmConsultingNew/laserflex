package laserflex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func AddCatalogDocumentElement(documentID int, productId int, quantity float64) error {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "catalog.document.element.add"

	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	requestBody := map[string]interface{}{
		"fields": map[string]interface{}{
			"docId":     documentID,
			"elementId": productId,
			"storeTo":   "1",           // ID склада
			"amount":    int(quantity), // Количество из массива
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %v", err)
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(responseData, &response); err != nil {
		return fmt.Errorf("error unmarshalling response: %v", err)
	}

	if _, ok := response["error"]; ok {
		log.Printf("Full response on error: %v", string(responseData))
		return fmt.Errorf("ошибка: %s", response["error_description"])
	}

	if _, ok := response["result"]; !ok {
		log.Printf("Unexpected response structure: %v", string(responseData))
		return fmt.Errorf("missing 'result' in response")
	}

	log.Printf("Товар добавлен в документ: ID документа: %s, ID товара: %d\n", documentID, productId)

	return nil
}
