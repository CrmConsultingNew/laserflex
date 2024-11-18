package deals

import (
	"bitrix_app/backend/bitrix/authorize"
	"log"
	"net/http"
)

var mockDealsJSON = `{"result":[{"ID":"43","TITLE":"Сделка #43","TYPE_ID":"SALE","STAGE_ID":"PREPAYMENT_INVOICE","PROBABILITY":null,"CURRENCY_ID":"RUB","OPPORTUNITY":"0.00","IS_MANUAL_OPPORTUNITY":"N","TAX_VALUE":"0.00","LEAD_ID":null,"COMPANY_ID":"0","CONTACT_ID":null,"QUOTE_ID":null,"BEGINDATE":"2024-03-20T03:00:00+03:00","CLOSEDATE":"2024-03-27T03:00:00+03:00","ASSIGNED_BY_ID":"15","CREATED_BY_ID":"15","MODIFY_BY_ID":"1","DATE_CREATE":"2024-03-20T21:33:17+03:00","DATE_MODIFY":"2024-03-30T00:42:52+03:00","OPENED":"Y","CLOSED":"N","COMMENTS":"","ADDITIONAL_INFO":null,"LOCATION_ID":null,"CATEGORY_ID":"0","STAGE_SEMANTIC_ID":"P","IS_NEW":"N","IS_RECURRING":"N","IS_RETURN_CUSTOMER":"N","IS_REPEATED_APPROACH":"N","SOURCE_ID":"","SOURCE_DESCRIPTION":null,"ORIGINATOR_ID":null,"ORIGIN_ID":null,"MOVED_BY_ID":"15","MOVED_TIME":"2024-03-20T21:33:45+03:00","LAST_ACTIVITY_TIME":"2024-03-21T21:33:45+03:00","UTM_SOURCE":null,"UTM_MEDIUM":null,"UTM_CAMPAIGN":null,"UTM_CONTENT":null,"UTM_TERM":null,"LAST_ACTIVITY_BY":"15"},{"ID":"65","TITLE":"Сделка #65","TYPE_ID":"SALE","STAGE_ID":"NEW","PROBABILITY":null,"CURRENCY_ID":"RUB","OPPORTUNITY":"0.00","IS_MANUAL_OPPORTUNITY":"N","TAX_VALUE":"0.00","LEAD_ID":null,"COMPANY_ID":"0","CONTACT_ID":null,"QUOTE_ID":null,"BEGINDATE":"2024-03-30T03:00:00+03:00","CLOSEDATE":"2024-04-06T03:00:00+03:00","ASSIGNED_BY_ID":"1","CREATED_BY_ID":"1","MODIFY_BY_ID":"1","DATE_CREATE":"2024-03-30T21:02:07+03:00","DATE_MODIFY":"2024-03-30T21:02:07+03:00","OPENED":"Y","CLOSED":"N","COMMENTS":"","ADDITIONAL_INFO":null,"LOCATION_ID":null,"CATEGORY_ID":"0","STAGE_SEMANTIC_ID":"P","IS_NEW":"Y","IS_RECURRING":"N","IS_RETURN_CUSTOMER":"N","IS_REPEATED_APPROACH":"N","SOURCE_ID":"","SOURCE_DESCRIPTION":null,"ORIGINATOR_ID":null,"ORIGIN_ID":null,"MOVED_BY_ID":"1","MOVED_TIME":"2024-03-30T21:02:07+03:00","LAST_ACTIVITY_TIME":"2024-03-30T21:02:07+03:00","UTM_SOURCE":null,"UTM_MEDIUM":null,"UTM_CAMPAIGN":null,"UTM_CONTENT":null,"UTM_TERM":null,"LAST_ACTIVITY_BY":"1"},{"ID":"67","TITLE":"Сделка #67","TYPE_ID":"SALE","STAGE_ID":"NEW","PROBABILITY":null,"CURRENCY_ID":"RUB","OPPORTUNITY":"0.00","IS_MANUAL_OPPORTUNITY":"N","TAX_VALUE":"0.00","LEAD_ID":null,"COMPANY_ID":"0","CONTACT_ID":null,"QUOTE_ID":null,"BEGINDATE":"2024-03-30T03:00:00+03:00","CLOSEDATE":"2024-04-06T03:00:00+03:00","ASSIGNED_BY_ID":"1","CREATED_BY_ID":"1","MODIFY_BY_ID":"1","DATE_CREATE":"2024-03-30T21:02:14+03:00","DATE_MODIFY":"2024-03-30T21:02:14+03:00","OPENED":"Y","CLOSED":"N","COMMENTS":"","ADDITIONAL_INFO":null,"LOCATION_ID":null,"CATEGORY_ID":"0","STAGE_SEMANTIC_ID":"P","IS_NEW":"Y","IS_RECURRING":"N","IS_RETURN_CUSTOMER":"N","IS_REPEATED_APPROACH":"N","SOURCE_ID":"","SOURCE_DESCRIPTION":null,"ORIGINATOR_ID":null,"ORIGIN_ID":null,"MOVED_BY_ID":"1","MOVED_TIME":"2024-03-30T21:02:14+03:00","LAST_ACTIVITY_TIME":"2024-03-30T21:02:14+03:00","UTM_SOURCE":null,"UTM_MEDIUM":null,"UTM_CAMPAIGN":null,"UTM_CONTENT":null,"UTM_TERM":null,"LAST_ACTIVITY_BY":"1"},{"ID":"69","TITLE":"Сделка #69","TYPE_ID":"SALE","STAGE_ID":"NEW","PROBABILITY":null,"CURRENCY_ID":"RUB","OPPORTUNITY":"0.00","IS_MANUAL_OPPORTUNITY":"N","TAX_VALUE":"0.00","LEAD_ID":null,"COMPANY_ID":"0","CONTACT_ID":null,"QUOTE_ID":null,"BEGINDATE":"2024-03-30T03:00:00+03:00","CLOSEDATE":"2024-04-06T03:00:00+03:00","ASSIGNED_BY_ID":"1","CREATED_BY_ID":"1","MODIFY_BY_ID":"1","DATE_CREATE":"2024-03-30T21:02:19+03:00","DATE_MODIFY":"2024-03-30T21:02:19+03:00","OPENED":"Y","CLOSED":"N","COMMENTS":"","ADDITIONAL_INFO":null,"LOCATION_ID":null,"CATEGORY_ID":"0","STAGE_SEMANTIC_ID":"P","IS_NEW":"Y","IS_RECURRING":"N","IS_RETURN_CUSTOMER":"N","IS_REPEATED_APPROACH":"N","SOURCE_ID":"","SOURCE_DESCRIPTION":null,"ORIGINATOR_ID":null,"ORIGIN_ID":null,"MOVED_BY_ID":"1","MOVED_TIME":"2024-03-30T21:02:19+03:00","LAST_ACTIVITY_TIME":"2024-03-30T21:02:19+03:00","UTM_SOURCE":null,"UTM_MEDIUM":null,"UTM_CAMPAIGN":null,"UTM_CONTENT":null,"UTM_TERM":null,"LAST_ACTIVITY_BY":"1"},{"ID":"71","TITLE":"Сделка #71","TYPE_ID":"SALE","STAGE_ID":"NEW","PROBABILITY":null,"CURRENCY_ID":"RUB","OPPORTUNITY":"0.00","IS_MANUAL_OPPORTUNITY":"N","TAX_VALUE":"0.00","LEAD_ID":null,"COMPANY_ID":"0","CONTACT_ID":null,"QUOTE_ID":null,"BEGINDATE":"2024-03-30T03:00:00+03:00","CLOSEDATE":"2024-04-06T03:00:00+03:00","ASSIGNED_BY_ID":"1","CREATED_BY_ID":"1","MODIFY_BY_ID":"1","DATE_CREATE":"2024-03-30T21:02:23+03:00","DATE_MODIFY":"2024-03-30T21:02:23+03:00","OPENED":"Y","CLOSED":"N","COMMENTS":"","ADDITIONAL_INFO":null,"LOCATION_ID":null,"CATEGORY_ID":"0","STAGE_SEMANTIC_ID":"P","IS_NEW":"Y","IS_RECURRING":"N","IS_RETURN_CUSTOMER":"N","IS_REPEATED_APPROACH":"N","SOURCE_ID":"","SOURCE_DESCRIPTION":null,"ORIGINATOR_ID":null,"ORIGIN_ID":null,"MOVED_BY_ID":"1","MOVED_TIME":"2024-03-30T21:02:23+03:00","LAST_ACTIVITY_TIME":"2024-03-30T21:02:23+03:00","UTM_SOURCE":null,"UTM_MEDIUM":null,"UTM_CAMPAIGN":null,"UTM_CONTENT":null,"UTM_TERM":null,"LAST_ACTIVITY_BY":"1"},{"ID":"73","TITLE":"Сделка #73","TYPE_ID":"SALE","STAGE_ID":"NEW","PROBABILITY":null,"CURRENCY_ID":"RUB","OPPORTUNITY":"0.00","IS_MANUAL_OPPORTUNITY":"N","TAX_VALUE":"0.00","LEAD_ID":null,"COMPANY_ID":"0","CONTACT_ID":null,"QUOTE_ID":null,"BEGINDATE":"2024-03-30T03:00:00+03:00","CLOSEDATE":"2024-04-06T03:00:00+03:00","ASSIGNED_BY_ID":"1","CREATED_BY_ID":"1","MODIFY_BY_ID":"1","DATE_CREATE":"2024-03-30T21:02:29+03:00","DATE_MODIFY":"2024-03-30T21:02:29+03:00","OPENED":"Y","CLOSED":"N","COMMENTS":"","ADDITIONAL_INFO":null,"LOCATION_ID":null,"CATEGORY_ID":"0","STAGE_SEMANTIC_ID":"P","IS_NEW":"Y","IS_RECURRING":"N","IS_RETURN_CUSTOMER":"N","IS_REPEATED_APPROACH":"N","SOURCE_ID":"","SOURCE_DESCRIPTION":null,"ORIGINATOR_ID":null,"ORIGIN_ID":null,"MOVED_BY_ID":"1","MOVED_TIME":"2024-03-30T21:02:29+03:00","LAST_ACTIVITY_TIME":"2024-03-30T21:02:29+03:00","UTM_SOURCE":null,"UTM_MEDIUM":null,"UTM_CAMPAIGN":null,"UTM_CONTENT":null,"UTM_TERM":null,"LAST_ACTIVITY_BY":"1"}],"total":6,"time":{"start":1711825725.99915,"finish":1711825726.0454381,"duration":0.046288013458251953,"processing":0.024226903915405273,"date_start":"2024-03-30T22:08:45+03:00","date_finish":"2024-03-30T22:08:46+03:00","operating_reset_at":1711826326,"operating":0}}
`

// Mock function to return test deals data
func GetDealsMock() ([]byte, error) {
	return []byte(mockDealsJSON), nil
}

func TransferDealsOnVueMock(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "GET" {
		log.Println("TransferDealsOnVue method is GET")
		deals, err := GetDealsMock() // Replace "your-auth-id" with your actual authorization GlobalAuthId
		if err != nil {
			log.Println("error getting service: ", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set content type to application/json
		w.Header().Set("Content-Type", "application/json")

		// Write the JSON data to the response
		_, writeErr := w.Write(deals)
		if writeErr != nil {
			log.Println("error writing service to response: ", writeErr.Error())
		}

	} else {
		log.Println("something wrong...")
	}
}

func TransferDealsOnVue(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "GET" {
		log.Println("TransferDealsOnVue method is GET")
		deals, err := GetDeals(authorize.GlobalAuthId)
		if err != nil {
			log.Println("error getting service: ", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set content type to application/json
		w.Header().Set("Content-Type", "application/json")

		// Log the deals for debugging
		//fmt.Println("TransferDealsOnVue deals: ", string(deals))

		// Write the JSON data to the response
		_, writeErr := w.Write(deals)
		if writeErr != nil {
			log.Println("error writing service to response: ", writeErr.Error())
		}

	} else {
		log.Println("something wrong...")
	}

}
