package laserflex

import (
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
	redirectURL := "https://crmconsulting-api.ru/"

	// Use http.Redirect to redirect the client
	// The http.StatusFound status code is commonly used for redirects
	http.Redirect(w, r, redirectURL, http.StatusFound)

	fmt.Println("redirect is done...")
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
