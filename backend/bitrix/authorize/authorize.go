package authorize

import (
	"bitrix_app/backend/bitrix/authorize/auth"
	"bitrix_app/backend/bitrix/service/leads"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var GlobalAuthId string

// Основная функция, которая запускает GetLeads и UpdateLeads
func ConnectionBitrixLocalAppNew(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection is starting...")
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading request body:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	log.Println("resp_at_new_connection:", string(bs))
	defer r.Body.Close()

	authValues := ParseValues(w, bs)
	fmt.Printf("authValues.AuthID : %s, authValues.MemberID: %s", authValues.AuthID, authValues.MemberID)

	// Бесконечный цикл для выполнения GetLeads и UpdateLeads
	for {
		// Получаем список лидов (20 запросов)
		getLeads, err := leads.GetLeads(authValues.AuthID)
		if err != nil {
			log.Println("Could not get leads:", err)
			continue // Продолжаем цикл даже при ошибке
		}

		// Если лиды получены, выполняем их обновление
		if len(getLeads) > 0 {
			for _, v := range getLeads {
				_, err := leads.UpdateLeads(v.ID, authValues.AuthID)
				if err != nil {
					log.Println("Could not update leads:", err)
				}
			}
		}

		// Лог после завершения цикла обновления
		log.Printf("Итого выполнено %d запросов в GetLeads и %d запросов в UpdateLeads", leads.GetLeadsCount, leads.UpdateLeadsCount)

		// Пауза перед началом следующего цикла
		time.Sleep(5 * time.Second) // Пауза, чтобы не перегружать API
	}
}

func ConnectionBitrixLocalApp(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection is starting...")
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading request body:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	log.Println("resp_at_first:", string(bs))
	defer r.Body.Close()

	authValues := ParseValues(w, bs) //todo here we must to add this data in dbase?
	fmt.Printf("authValues.AuthID : %s, authValues.MemberID: %s", authValues.AuthID, authValues.MemberID)

	//w.Write([]byte(authValues.AuthID))
	redirectURL := "https://crmconsulting-api.ru/"

	// Use http.Redirect to redirect the client
	// The http.StatusFound status code is commonly used for redirects
	http.Redirect(w, r, redirectURL, http.StatusFound)

	fmt.Println("redirect is done...")
	GlobalAuthId = authValues.AuthID

	//events.OnCrmDealAddEventRegistration(authValues.AuthID) //todo return this method
}

func ParseValues(w http.ResponseWriter, bs []byte) auth.Request {
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

	authValues := auth.Request{
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
