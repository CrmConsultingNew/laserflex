package events

import (
	"bitrix_app/backend/bitrix/endpoints"
	"fmt"
	"io"
	"log"
	"net/http"
)

func OnCrmDealAddEventRegistration(authId string) {

	method := "POST"
	event := "OnCrmDealAdd"
	handler := fmt.Sprintf("%s/api/event_deal_add", endpoints.BitrixDomain)

	// Format the URL with the provided authID parameter
	requestUrl := fmt.Sprintf("%s/rest/event.bind.json?auth=%s&auth_type=0&event=%s&event_type=offline&handler=%s", endpoints.BitrixDomain, authId, event, handler)

	req, err := http.NewRequest(method, requestUrl, nil)
	if err != nil {
		log.Println("error creating new request:", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error sending request:", err)
	}
	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response body:", err)
	}
	log.Println("resp_at_last_AddDeal:", string(bz))

}
