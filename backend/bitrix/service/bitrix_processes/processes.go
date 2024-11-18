package bitrix_processes

import (
	"bitrix_app/backend/bitrix/authorize"
	"bitrix_app/backend/bitrix/endpoints"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// ProcessesRequestBody represents the request body structure for the Bitrix API request.
type ProcessesRequestBody struct {
	IBLOCKTYPEID string `json:"IBLOCK_TYPE_ID"`
	IBLOCKID     string `json:"IBLOCK_ID"`
	Filter       struct {
		PROPERTY632 map[string]string `json:"PROPERTY_632"`
		PROPERTY634 map[string]string `json:"PROPERTY_634"`
	} `json:"filter"`
}

// ProcessesResponse - ДДС (Бизнес-процесс в ленте новостей)
// ProcessesResponse represents the structure of the response that will be sent to the frontend.
type ProcessesResponse struct {
	ID          string `json:"ID"`
	Property628 string `json:"PROPERTY_628"` // Сумма
	Property638 string `json:"PROPERTY_638"` // Дата операции - Старое поле  PROPERTY_732
}

func GetProcessesListHandler(w http.ResponseWriter, r *http.Request) {
	companyID := r.URL.Query().Get("id")
	if companyID == "" {
		http.Error(w, "Company ID is missing", http.StatusBadRequest)
		return
	}

	property634Value := "788" // замените на реальное значение property634

	// Вызов функции для получения списка процессов
	processes, err := GetProcessesList(authorize.GlobalAuthId, companyID, property634Value)
	if err != nil {
		http.Error(w, "Error fetching processes data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Преобразование данных в JSON для отправки на фронтенд
	responseData, err := json.Marshal(processes)
	if err != nil {
		http.Error(w, "Error marshaling processes data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Установка заголовков и отправка данных на фронтенд
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

// GetProcessesList sends a request to the Bitrix API and returns the response data.
// GetProcessesList sends a request to the Bitrix API and returns the response data.
func GetProcessesList(authID, companyID, property634Value string) ([]ProcessesResponse, error) {
	log.Println("GetProcessesList was started:")
	bitrixMethod := "lists.element.get"
	requestURL := fmt.Sprintf("%s/rest/%s?auth=%s", endpoints.BitrixDomain, bitrixMethod, authID)

	// Construct the new request body with the provided company ID and PROPERTY_634 value
	body := ProcessesRequestBody{
		IBLOCKTYPEID: "bitrix_processes",
		IBLOCKID:     "108",
		Filter: struct {
			PROPERTY632 map[string]string `json:"PROPERTY_632"` // поступления от кого - Контрагент
			PROPERTY634 map[string]string `json:"PROPERTY_634"` // приход / расход
		}{
			PROPERTY632: map[string]string{
				"28644": fmt.Sprintf("CO_%s", companyID), // Формируем фильтр по ID компании
			},
			PROPERTY634: map[string]string{
				"38158": property634Value, // Используем значение, переданное в функцию
			},
		},
	}

	// Marshal the request body into JSON
	jsonData, err := json.Marshal(body)
	if err != nil {
		log.Println("Error marshaling request body:", err)
		return nil, err
	}

	// Create a new request with JSON body
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating new requests:", err)
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
	log.Println("GetProcessesList Response:", string(responseData))

	// Parse the response into a slice of processes
	var rawResponse struct {
		Result []map[string]interface{} `json:"result"`
	}
	if err := json.Unmarshal(responseData, &rawResponse); err != nil {
		log.Println("Error unmarshaling response:", err)
		return nil, err
	}

	// Extract the needed values from PROPERTY_628 and Property638
	var processes []ProcessesResponse
	for _, item := range rawResponse.Result {
		var property628, Property638 string

		if property, ok := item["PROPERTY_628"].(map[string]interface{}); ok {
			for _, v := range property {
				if value, ok := v.(string); ok {
					property628 = value
					break // Берем только первое значение
				}
			}
		}

		if property, ok := item["PROPERTY_638"].(map[string]interface{}); ok {
			for _, v := range property {
				if value, ok := v.(string); ok {
					Property638 = value
					break // Берем только первое значение
				}
			}
		}

		processes = append(processes, ProcessesResponse{
			ID:          item["ID"].(string),
			Property628: property628,
			Property638: Property638,
		})
	}
	log.Println("GetProcessesList was ended:")
	return processes, nil
}

// ProcessesHandler формирует и возвращает документ .docx
/*func ProcessesHandler(w http.ResponseWriter, r *http.Request) {
	companyID := r.URL.Query().Get("id")
	if companyID == "" {
		http.Error(w, "Company ID is missing", http.StatusBadRequest)
		return
	}

	log.Println("ProcessesHandler: Starting document generation...")

	var tableDataLeft [][]string

	// Получаем данные из Bitrix и добавляем их в таблицу
	processesData, err := GetProcessesList(authorize.GlobalAuthId, companyID)
	if err != nil {
		http.Error(w, "Error fetching processes data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	AddProcessesToTable(&tableDataLeft, processesData)

	// Создание документа Word
	err = office.StartWord(
		"01.01.2024",       // dateFrom
		"31.12.2024",       // dateTo
		"15.01.2024",       // dateForInvoice
		"20.02.2024",       // dateForReturnPayment
		"12345",            // numberOfCompletedWorks
		"10.03.2024",       // dateOfCompletedWorks
		"31.12.2024",       // dateForSaldo
		"10000",            // sumOfSaldo
		"ООО «Компания 2»", // secondCompany
		tableDataLeft,
	)
	if err != nil {
		http.Error(w, "Error creating Word document: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Путь к сгенерированному файлу
	filePath := "tables.docx"

	// Открытие файла для чтения
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Ошибка при открытии файла: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Чтение файла в память
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		http.Error(w, "Ошибка при чтении файла: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Настройка заголовков для скачивания файла
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))

	// Отправка файла клиенту
	_, err = w.Write(buf.Bytes())
	if err != nil {
		log.Println("Write error ")
		return
	}

	log.Println("ProcessesHandler: Document successfully sent to client.")
}*/
