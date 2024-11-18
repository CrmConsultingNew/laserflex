package routes

import (
	"bitrix_app/backend/bitrix/authorize"
	"bitrix_app/backend/bitrix/service/bitrix_processes"
	"bitrix_app/backend/bitrix/service/comments"
	"bitrix_app/backend/bitrix/service/companies"
	"bitrix_app/backend/bitrix/service/deals"
	"bitrix_app/backend/bitrix/service/description"
	"bitrix_app/backend/bitrix/service/docs"
	"bitrix_app/backend/bitrix/service/events"
	"bitrix_app/backend/bitrix/service/settings"
	smart_processes "bitrix_app/backend/bitrix/service/smart-processes"
	"bitrix_app/backend/bitrix/test"
	"bitrix_app/backend/laserflex"
	"bitrix_app/backend/widget"
	"net/http"
)

func Router() {

	http.HandleFunc("/send_file", laserflex.LaserflexGetFile)

	////////////////////////////////////////////////////////////////

	http.HandleFunc("/api/connect_widget", widget.ConnectionBitrixWidget)

	http.HandleFunc("/api/connect_widget", widget.ConnectionBitrixWidget)
	http.HandleFunc("/api/widget_data", widget.SendDataForWidgetForm)
	http.HandleFunc("/api/form_data", widget.GetDataFromWidgetForm)

	http.HandleFunc("/api/connect", authorize.ConnectionBitrixLocalApp)
	http.HandleFunc("/api/companies", companies.CompaniesHandler)
	http.HandleFunc("/api/company", companies.CompanyHandler)
	http.HandleFunc("/api/processes", bitrix_processes.GetProcessesListHandler)

	http.HandleFunc("/api/items", smart_processes.GetItemsByCompanyHandler)

	http.HandleFunc("/api/deals_get", deals.TransferDealsOnVue)
	http.HandleFunc("/api/event_deal_add", events.OnCrmDealAddEvent)

	http.HandleFunc("/api/documents/", docs.DocumentHandler)
	http.HandleFunc("/api/comments/", comments.CommentsHandler)
	http.HandleFunc("/api/description/", description.DescriptionHandler)

	http.HandleFunc("/api/save_settings", settings.SaveSettingsHandler)

	//http.HandleFunc("/api/gpt", chatgpt.SendRequest)

	http.HandleFunc("/api/user-redirect/", test.UserRedirect)
	http.HandleFunc("/api/user-form", test.UserForm)
	http.HandleFunc("/api/deal_id", test.GetWebhookWithDealId)

	http.HandleFunc("/api/sended_sms", test.SendedSms)
	http.HandleFunc("/api/sended_done_sms", test.SendedDoneSms)

	//http.HandleFunc("/api/check_widget", widget.CheckWidget) //here we create widget in bitrix

	/*c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Change this to the specific origin of your Vue.js app in a production environment.
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	http.Handle("/api/auth_page", c.Handler(http.HandlerFunc(repository.AuthPage)))
	http.Handle("/api/login_page", c.Handler(http.HandlerFunc(repository.LoginPage)))
	http.HandleFunc("/api/redirect", repository.RedirectPage) //here user redirects from login page*/

	http.HandleFunc("/api/redirect", deals.ConnectionBitrixLogger)

}
