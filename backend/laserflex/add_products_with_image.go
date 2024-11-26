package laserflex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ProductCreationResponse struct {
	Result int `json:"result"`
}

func AddProductsWithImage(product Product, sectionID string) (int, error) {
	webHookUrl := "https://bitrix.laser-flex.ru/rest/149/5cycej8804ip47im/"
	bitrixMethod := "crm.product.add"

	requestURL := fmt.Sprintf("%s%s", webHookUrl, bitrixMethod)

	requestBody := map[string]interface{}{
		"fields": map[string]interface{}{
			"NAME":       product.Name,
			"SECTION_ID": sectionID,
			"PROPERTY_197": []map[string]interface{}{
				{
					"valueId":  0,
					"fileData": []string{"avatar.png", product.ImageBase64},
				},
			},
			"PROPERTY_198": product.Material,
			"PROPERTY_199": product.Laser,
			"PROPERTY_200": product.Bend,
			"PROPERTY_201": product.Paint,
			"PROPERTY_202": product.Production,
			"PROPERTY_203": product.PipeCutting,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return 0, fmt.Errorf("error marshalling request body: %v", err)
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	var response ProductCreationResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return 0, fmt.Errorf("error unmarshalling response: %v", err)
	}

	if response.Result == 0 {
		return 0, fmt.Errorf("failed to create product, response: %s", string(responseData))
	}

	log.Println("Product added with ID:", response.Result)
	return response.Result, nil
}
