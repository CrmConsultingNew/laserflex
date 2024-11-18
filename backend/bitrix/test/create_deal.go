package test

import (
	"bitrix_app/backend/bitrix/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func CreateDeal(commentary string, category string, link string, contactId string, branch string, rating int, dateCreate string, visitDate string, stageId string) error {
	method := "POST"
	requestBody := fmt.Sprintf(`{"fields":
{"TITLE":"Сбор обратной связи harizma-service",
"UF_CRM_1712936072":"%s", 
"CATEGORY_ID":"%s",
"UF_CRM_1712927864":"%s",
"CONTACT_ID":"%s",
"DATE_CREATE":"%s",
"UF_CRM_1690982742603":"%v",
"UF_CRM_1712927909": "%v",
"UF_CRM_1690209734961": "%v",
"STAGE_ID": "%s"
}
}`, commentary, category, link, contactId, dateCreate, branch, rating, visitDate, stageId)

	//UF_CRM_1712927864 = Ссылка на сделку (NPS)
	//UF_CRM_1690982742603 = Адрес филиала
	//UF_CRM_1712927909 = Оценка (NPS)
	//UF_CRM_1712936072 = Комментарий (NPS)
	//UF_CRM_1690209734961 = Дата и время визита (для фильтра)
	// Convert the JSON string to a byte slice
	body := []byte(requestBody)

	webHookUrl := os.Getenv("WEBHOOK_URL_HARIZMA_DEAL_ADD")
	// Format the URL with the provided authID parameter

	req, err := http.NewRequest(method, webHookUrl, bytes.NewReader(body))
	if err != nil {
		log.Println("error creating new request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json") // Set the content type to application/json

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response body:", err)
		return err
	}

	log.Println("resp_at_last_AddDeal:", string(bz))
	return nil
}

func GetDealById(id string) (models.ApiResponseHarizma, error) {
	//contactId, branch, dateCreate
	method := "POST"
	webhookUrl := os.Getenv("WEBHOOK_URL_HARIZMA_DEAL_GET")
	requestBody := fmt.Sprintf(`{
"id":"%s"
}`, id)

	// Convert the JSON string to a byte slice
	body := []byte(requestBody)

	req, err := http.NewRequest(method, webhookUrl, bytes.NewReader(body))
	if err != nil {
		log.Println("error creating new request:", err)
		return models.ApiResponseHarizma{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error sending request:", err)
		return models.ApiResponseHarizma{}, err
	}
	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response body:", err)
		return models.ApiResponseHarizma{}, err
	}

	log.Println("resp_at_last_AddDeal:", string(bz))

	var apiResponse models.ApiResponseHarizma
	if err := json.Unmarshal(bz, &apiResponse); err != nil {
		log.Printf("error unmarshalling response to ApiResponse: %v", err)
		return models.ApiResponseHarizma{}, err
	}

	return apiResponse, nil
}
