package leads

import (
	"bitrix_app/backend/bitrix/endpoints"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func AddLead(titleID string, authID string) ([]Leads, error) {
	method := "POST" // Используем POST для отправки данных с телом

	// Форматируем URL с использованием authID
	requestUrl := fmt.Sprintf("%s/rest/crm.lead.add?auth=%s", endpoints.NewBitrixDomain, authID)

	// Тело запроса, подставляем значение titleID
	body := fmt.Sprintf(`{
		"TITLE": "TEST%s"
	}`, titleID)

	// Создаем новый HTTP-запрос с телом
	req, err := http.NewRequest(method, requestUrl, strings.NewReader(body))
	if err != nil {
		log.Println("Ошибка при создании запроса:", err)
		return nil, err
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Ошибка при отправке запроса:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Читаем ответ от сервера
	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Ошибка при чтении тела ответа:", err)
		return nil, err
	}

	log.Println("Ответ от Bitrix (leads):", string(bz))

	// Парсим ответ в структуру
	var apiResponse ApiResponseLeads
	if err := json.Unmarshal(bz, &apiResponse); err != nil {
		log.Printf("Ошибка при Unmarshal ответа в структуру: %v", err)
		return nil, err
	}

	// Возвращаем массив лидов
	return apiResponse.Result, nil
}

// Обработка ошибки и повтор запроса
func retryRequest(req *http.Request) (*http.Response, error) {
	for {
		// Копируем запрос для корректного повторения
		reqCopy := req.Clone(req.Context())

		resp, err := http.DefaultClient.Do(reqCopy)
		if err != nil {
			return nil, err
		}

		// Проверяем статус-код 503, чтобы сделать паузу и повторить запрос
		if resp.StatusCode == 503 {
			log.Println("Достигнут лимит запросов (503), ожидаем перед повтором...")
			resp.Body.Close() // Закрываем тело перед повтором
			time.Sleep(2 * time.Minute)
			continue // Переходим к следующей попытке
		}

		// Проверяем, если статус успешный (менее 300)
		if resp.StatusCode < 300 {
			return resp, nil
		}

		// Закрываем тело, если есть другие ошибки
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

var GetLeadsCount int    // Счётчик запросов для GetLeads
var UpdateLeadsCount int // Счётчик запросов для UpdateLeads

func GetLeads(authID string) ([]Leads, error) {
	method := "POST"
	allLeads := []Leads{}
	start := 0
	limit := 50 // Ограничение в 50 записей
	requestCount := 0

	// Определяем лимит по дате
	//dateLimit := "2024-09-30T00:00:00+03:00"

	for i := 0; i < 20; i++ { // Цикл для выполнения 20 запросов
		// Формируем URL и тело запроса
		requestUrl := fmt.Sprintf("%s/rest/crm.lead.list?auth=%s", endpoints.NewBitrixDomain, authID)
		body := fmt.Sprintf(`{
			"filter": {
				"STATUS_ID": "IN_PROCESS",
                "!COMMENTS":"date_after_1_october_2024"
			},
			"start": %d,
			"limit": %d
		}`, start, limit)

		// Создаем новый HTTP-запрос
		req, err := http.NewRequest(method, requestUrl, strings.NewReader(body))
		if err != nil {
			log.Println("Ошибка при создании запроса:", err)
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		// Используем retryRequest для обработки лимитов и повторов
		resp, err := retryRequest(req)
		if err != nil {
			log.Println("Ошибка при отправке запроса:", err)
			return nil, err
		}
		defer resp.Body.Close()

		// Читаем тело ответа
		bz, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Ошибка при чтении тела ответа:", err)
			return nil, err
		}

		log.Println("Ответ от Bitrix (leads):", string(bz))

		// Обрабатываем JSON-ответ
		var apiResponse struct {
			Result []Leads `json:"result"`
			Next   int     `json:"next"`
		}
		if err := json.Unmarshal(bz, &apiResponse); err != nil {
			log.Printf("Ошибка при Unmarshal ответа в структуру: %v", err)
			return nil, err
		}

		allLeads = append(allLeads, apiResponse.Result...)
		GetLeadsCount++
		log.Printf("Запрос № %d выполнен в GetLeads, получено лидов: %d", GetLeadsCount, len(apiResponse.Result))

		// Если нет следующей страницы, завершаем цикл
		if apiResponse.Next == 0 {
			break
		}

		start = apiResponse.Next
		requestCount++

		// Пауза между запросами
		time.Sleep(2 * time.Second)
	}

	return allLeads, nil
}

// Обновлённая функция UpdateLeads
func UpdateLeads(leadId string, authID string) ([]Leads, error) {
	method := "POST"
	allLeads := []Leads{}
	start := 0

	for {
		requestUrl := fmt.Sprintf("%s/rest/crm.lead.update?auth=%s", endpoints.NewBitrixDomain, authID)
		body := fmt.Sprintf(`{
			"id": "%s",
			"fields": {
				"STATUS_ID": "UC_PR7OP6",
                "UF_CRM_1728320843": "date_after_1_october_2024(last)"
			},
			"start": %d
		}`, leadId, start)

		// Создаем новый HTTP-запрос
		req, err := http.NewRequest(method, requestUrl, strings.NewReader(body))
		if err != nil {
			log.Println("Ошибка при создании запроса:", err)
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		// Используем retryRequest для обработки лимитов и повторов
		resp, err := retryRequest(req)
		if err != nil {
			log.Println("Ошибка при отправке запроса:", err)
			return nil, err
		}
		defer resp.Body.Close()

		// Читаем тело ответа
		bz, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Ошибка при чтении тела ответа:", err)
			return nil, err
		}

		log.Println("Ответ от Bitrix (leads_update):", string(bz))

		// Обрабатываем JSON-ответ
		var apiResponse struct {
			Result []Leads `json:"result"`
			Next   int     `json:"next"`
		}

		if err := json.Unmarshal(bz, &apiResponse); err != nil {
			log.Printf("Ошибка при Unmarshal ответа в структуру: %v", err)
			return nil, err
		}

		allLeads = append(allLeads, apiResponse.Result...)
		UpdateLeadsCount++
		log.Printf("Запрос № %d выполнен в UpdateLeads для лида %s", UpdateLeadsCount, leadId)

		if apiResponse.Next == 0 {
			break
		}

		start = apiResponse.Next

		// Пауза между запросами
		time.Sleep(2 * time.Second)
	}

	return allLeads, nil
}
