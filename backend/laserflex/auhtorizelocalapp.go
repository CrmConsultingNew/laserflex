package laserflex

import (
	"bitrix_app/backend/bitrix/endpoints"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Request struct {
	AuthID           string `json:"auth_id"`
	AuthExpires      int    `json:"auth_expires"`
	RefreshID        string `json:"refresh_id"`
	MemberID         string `json:"member_id"`
	Status           string `json:"status"`
	Placement        string `json:"placement"`
	PlacementOptions string `json:"placement_options"`
}

var BlobalAuthIdLaserflex string

func AuthorizeEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection is starting...")
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading request body:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	log.Println("resp_at_first:", string(bs))
	defer r.Body.Close()

	authValues := ParseValuesLaserflex(w, bs) //todo here we must to add this data in dbase?
	fmt.Printf("authValues.AuthID : %s, authValues.MemberID: %s", authValues.AuthID, authValues.MemberID)

	//w.Write([]byte(authValues.AuthID))
	//redirectURL := "https://bitrix.laser-flex.ru/marketplace/app/25/"

	// Use http.Redirect to redirect the client
	// The http.StatusFound status code is commonly used for redirects
	//http.Redirect(w, r, redirectURL, http.StatusFound)

	//fmt.Println("redirect is done...")
	BlobalAuthIdLaserflex = authValues.AuthID
}

func ParseValuesLaserflex(w http.ResponseWriter, bs []byte) Request {
	values, err := url.ParseQuery(string(bs))
	if err != nil {
		log.Println("error parsing query:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	authExpires, err := strconv.Atoi(values.Get("AUTH_EXPIRES"))
	if err != nil {
		log.Println("error converting AUTH_EXPIRES to int:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	authValues := Request{
		AuthID:           values.Get("AUTH_ID"),
		AuthExpires:      authExpires,
		RefreshID:        values.Get("REFRESH_ID"),
		MemberID:         values.Get("member_id"),
		Status:           values.Get("status"),
		Placement:        values.Get("PLACEMENT"),
		PlacementOptions: values.Get("PLACEMENT_OPTIONS"),
	}
	return authValues
}

func GetAllCompaniesList(authKey string) ([]Company, error) {
	bitrixMethod := "crm.company.list"
	allCompanies := []Company{}
	start := 0

	for {
		requestURL := fmt.Sprintf("%s/rest/%s?auth=%s", endpoints.BitrixDomain, bitrixMethod, authKey)

		// Prepare request body with pagination
		requestBody := map[string]interface{}{
			"start": start, // Начало выборки
		}

		// Marshal the request body into JSON
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			log.Println("Error marshaling request body:", err)
			return nil, err
		}

		// Create a new request with JSON body
		req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("Error creating new request:", err)
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")

		// Send the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Error sending request:", err)
			return nil, err
		}
		defer resp.Body.Close()

		// Read and parse the response body
		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
			return nil, err
		}

		// Log the raw response for debugging
		//log.Println("GetCompaniesList Response:", string(responseData))

		// Parse the response into a slice of companies
		var response struct {
			Result []Company `json:"result"`
			Next   int       `json:"next"`
		}
		if err := json.Unmarshal(responseData, &response); err != nil {
			log.Println("Error unmarshaling response:", err)
			return nil, err
		}

		allCompanies = append(allCompanies, response.Result...)

		// Если поле next пустое или равно 0, завершить цикл
		if response.Next == 0 {
			break
		}

		start = response.Next
	}

	return allCompanies, nil
}

type Company struct {
	ID       string `json:"ID"`
	Title    string `json:"TITLE"`
	HasEmail string `json:"HAS_EMAIL"`
}
