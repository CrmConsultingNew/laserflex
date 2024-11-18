package widget

import (
	"bitrix_app/backend/bitrix/service/companies"
	"bitrix_app/backend/bitrix/service/deals"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

func stringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func GetDataFromWidgetForm(w http.ResponseWriter, r *http.Request) {
	// Чтение данных формы
	rdr, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Ошибка при чтении данных формы")
		return
	}
	log.Println("Cleared data from frontend: ", string(rdr))

	//log.Println("Request: ", string(rdr))
	// Преобразование полученных данных в структуру
	var formData map[string]string
	err = json.Unmarshal(rdr, &formData)
	if err != nil {
		log.Println("Ошибка при парсинге данных формы:", err)
		return
	}

	// Получение списка компаний
	list, err := companies.GetAllCompaniesList(GlobalAuthIdWidget)
	if err != nil {
		log.Println("Ошибка при получении списка компаний")
		return
	}

	// Проверяем, есть ли NaimenovanieKompanii в списке компаний
	for _, company := range list {
		if company.Title == formData["NaimenovanieKompanii"] {
			log.Println("Компания найдена, ID компании: ", company.ID)

			// Подготовка данных для обновления
			updatedCompany := companies.Company{
				Title:                            formData["NaimenovanieKompanii"],
				INN:                              stringToInt(formData["INN"]),
				NISHA:                            formData["Nisha"],
				GEOGRAPHIA:                       formData["Geographia"],
				KONECHNIY_PRODUCT:                formData["KonechniyProduct"],
				OBOROT_KOMPANII:                  stringToInt(formData["OborotKompanii"]),
				SREDNEMESYACHNAYA_VIRYCHKA:       stringToInt(formData["SrednemesyachnayaViryuchka"]),
				KOLICHESTVO_SOTRUDNIKOV:          stringToInt(formData["KolichestvoSotrudnikov"]),
				KOLICHESTVO_SOTRUDNIKOV_OP:       stringToInt(formData["KolichestvoSotrudnikovOP"]),
				EST_ROP:                          formData["EstROP"],
				EST_HR:                           formData["EstHR"],
				OSNOVNIE_KANALI_PRODAJ:           formData["OsnovnieKanaliProdazh"],
				SEZONNOST:                        formData["Sezonnost"],
				ZAPROS:                           formData["Zapros"],
				KAKIE_CELI_PERED_KOMPANIEI:       formData["CeliDoKontsaGoda"],
				CHTO_OJIDAETE_OT_SOTRUDNICHESTVA: formData["OzhidaniyaOtSotrudnichestva"],
			}

			// Обновляем компанию
			request, err := companies.UpdateCompany(company.ID, updatedCompany, GlobalAuthIdWidget)
			log.Println("request of update company:", request)
			if err != nil {
				log.Println("Ошибка при обновлении компании:", err)
			}
			// Подготовка данных для обновления сделки
			updatedDeal := deals.DealInfo{
				ID:                               DealIdGlobal, // Предполагается, что DealID передается в форме
				Title:                            formData["DealTitle"],
				Opportunity:                      formData["Opportunity"],
				CompanyID:                        stringToInt(company.ID),
				NISHA:                            formData["Nisha"],
				GEOGRAPHIA:                       formData["Geographia"],
				KONECHNIY_PRODUCT:                formData["KonechniyProduct"],
				KOLICHESTVO_SOTRUDNIKOV:          stringToInt(formData["KolichestvoSotrudnikov"]),
				KOLICHESTVO_SOTRUDNIKOV_OP:       stringToInt(formData["KolichestvoSotrudnikovOP"]),
				EST_ROP:                          formData["EstROP"],
				EST_HR:                           formData["EstHR"],
				OSNOVNIE_KANALI_PRODAJ:           formData["OsnovnieKanaliProdazh"],
				SEZONNOST:                        formData["Sezonnost"],
				ZAPROS:                           formData["Zapros"],
				KAKIE_CELI_PERED_KOMPANIEI:       formData["CeliDoKontsaGoda"],
				CHTO_OJIDAETE_OT_SOTRUDNICHESTVA: formData["OzhidaniyaOtSotrudnichestva"],
			}

			log.Println("updateDeal: ", updatedDeal)
			// Обновляем сделку
			_, err = deals.UpdateDeal(DealIdGlobal, updatedDeal, GlobalAuthIdWidget)
			if err != nil {
				log.Println("Ошибка при обновлении сделки:", err)
				return
			}
			return
		}
	}

	// Если компания не найдена, создаем новую
	newCompany := companies.Company{
		Title:                            formData["NaimenovanieKompanii"],
		INN:                              stringToInt(formData["INN"]),
		NISHA:                            formData["Nisha"],
		GEOGRAPHIA:                       formData["Geographia"],
		KONECHNIY_PRODUCT:                formData["KonechniyProduct"],
		OBOROT_KOMPANII:                  stringToInt(formData["OborotKompanii"]),
		SREDNEMESYACHNAYA_VIRYCHKA:       stringToInt(formData["SrednemesyachnayaViryuchka"]),
		KOLICHESTVO_SOTRUDNIKOV:          stringToInt(formData["KolichestvoSotrudnikov"]),
		KOLICHESTVO_SOTRUDNIKOV_OP:       stringToInt(formData["KolichestvoSotrudnikovOP"]),
		EST_ROP:                          formData["EstROP"],
		EST_HR:                           formData["EstHR"],
		OSNOVNIE_KANALI_PRODAJ:           formData["OsnovnieKanaliProdazh"],
		SEZONNOST:                        formData["Sezonnost"],
		ZAPROS:                           formData["Zapros"],
		KAKIE_CELI_PERED_KOMPANIEI:       formData["CeliDoKontsaGoda"],
		CHTO_OJIDAETE_OT_SOTRUDNICHESTVA: formData["OzhidaniyaOtSotrudnichestva"],
	}

	// Добавляем новую компанию
	addedCompany, err := companies.AddCompany(newCompany, GlobalAuthIdWidget)
	if err != nil {
		log.Println("Ошибка при добавлении компании:", err)
	}
	log.Println("Добавлена новая компания с ID: ", addedCompany)

	updatedDeal := deals.DealInfo{
		ID:                               DealIdGlobal, // Предполагается, что DealID передается в форме
		Title:                            formData["DealTitle"],
		Opportunity:                      formData["Opportunity"],
		CompanyID:                        addedCompany,
		NISHA:                            formData["Nisha"],
		GEOGRAPHIA:                       formData["Geographia"],
		KONECHNIY_PRODUCT:                formData["KonechniyProduct"],
		KOLICHESTVO_SOTRUDNIKOV:          stringToInt(formData["KolichestvoSotrudnikov"]),
		KOLICHESTVO_SOTRUDNIKOV_OP:       stringToInt(formData["KolichestvoSotrudnikovOP"]),
		EST_ROP:                          formData["EstROP"],
		EST_HR:                           formData["EstHR"],
		OSNOVNIE_KANALI_PRODAJ:           formData["OsnovnieKanaliProdazh"],
		SEZONNOST:                        formData["Sezonnost"],
		ZAPROS:                           formData["Zapros"],
		KAKIE_CELI_PERED_KOMPANIEI:       formData["CeliDoKontsaGoda"],
		CHTO_OJIDAETE_OT_SOTRUDNICHESTVA: formData["OzhidaniyaOtSotrudnichestva"],
	}

	// Преобразование updatedDeal в JSON для вывода в консоль
	updatedDealJSON, err := json.Marshal(updatedDeal)
	if err != nil {
		log.Println("Ошибка при преобразовании updatedDeal в JSON:", err)
	} else {
		log.Println("updatedDealSecond: ", string(updatedDealJSON))
	}

	// Обновляем сделку
	_, err = deals.UpdateDeal(DealIdGlobal, updatedDeal, GlobalAuthIdWidget)
	if err != nil {
		log.Println("Ошибка при обновлении сделки:", err)
		return
	}

}

func SendDataForWidgetForm(w http.ResponseWriter, r *http.Request) {

	/*log.Println("SendDataForWidgetForm DealIDGLOBAL: ", DealIdGlobal)

	// Получаем информацию о сделке
	dealInfo, err := deals.GetInfoAboutDealByID(DealIdGlobal, GlobalAuthIdWidget)
	if err != nil {
		log.Println("Error fetching deal info:", err)
		http.Error(w, "Error fetching deal info", http.StatusInternalServerError)
		return
	}*/

	// Если CompanyID отсутствует, получаем список компаний и отправляем всю структуру компании
	//if dealInfo.CompanyID == "0" {
	log.Println("Fetching list of companies")
	companiesList, err := companies.GetAllCompaniesList(GlobalAuthIdWidget)
	if err != nil {
		log.Println("Error fetching companies list:", err)
		http.Error(w, "Error fetching companies list", http.StatusInternalServerError)
	}

	log.Println("companiesList", companiesList)
	// Отправляем список компаний на фронтенд
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(companiesList)
	return
	//}

	/*// Если есть CompanyID, получаем информацию о компании
	if dealInfo.CompanyID != "0" {
		companyInfo, err := companies.GetCompanyByID(dealInfo.CompanyID, GlobalAuthIdWidget)
		if err != nil {
			log.Println("Error fetching company info:", err)
			http.Error(w, "Error fetching company info", http.StatusInternalServerError)
			return
		}

		// Преобразуем и отправляем данные компании
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(companyInfo)
	}*/
}
