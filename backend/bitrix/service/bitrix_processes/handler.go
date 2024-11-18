package bitrix_processes

import (
	"bitrix_app/backend/bitrix/authorize"
	"bitrix_app/backend/bitrix/service/requisites"
	smart_processes "bitrix_app/backend/bitrix/service/smart-processes"
	"bitrix_app/backend/office"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

// Структура для объединения данных оплаты и актов
type DocumentEntry struct {
	Type        string    // Тип записи: "оплата" или "акт"
	Date        time.Time // Дата операции
	Description string    // Описание записи
	Credit      string    // Кредит (для оплат)
	Debit       string    // Дебет (для актов)
}

// AddProcessesToTable добавляет данные оплаты в общий список
func AddProcessesToTable(entries *[]DocumentEntry, processes []ProcessesResponse) {
	for _, process := range processes {
		// Парсим дату
		processDate, err := time.Parse("02.01.2006", process.Property638)
		if err == nil {
			// Добавляем запись об оплате в список
			*entries = append(*entries, DocumentEntry{
				Type:        "оплата",
				Date:        processDate,
				Description: fmt.Sprintf("Оплата по счету от %s", processDate.Format("02.01.2006")),
				Credit:      process.Property628, // Кредит для оплат
			})
		} else {
			log.Printf("Error formatting process date for process ID: %v, error: %v", process.ID, err)
		}
	}
}

// AddItemsToTable добавляет данные актов в общий список
func AddItemsToTable(entries *[]DocumentEntry, items []smart_processes.ItemsResponse) {
	for _, item := range items {
		// Форматируем дату акта
		itemDate, err := time.Parse("02.01.2006", item.UFCrm1712128088)
		if err == nil {
			// Преобразуем Opportunity из float64 в строку
			opportunityStr := strconv.FormatFloat(item.Opportunity, 'f', 0, 64)

			// Добавляем запись акта в список
			*entries = append(*entries, DocumentEntry{
				Type:        "акт",
				Date:        itemDate,
				Description: fmt.Sprintf("Акт выполненных работ от %s", itemDate.Format("02.01.2006")),
				Debit:       opportunityStr, // Дебет для актов
			})
		} else {
			log.Printf("Error formatting item date for item ID: %v, error: %v", item.ID, err)
		}
	}
}

// FileGeneratorHandler обработчик для генерации файлов
func FileGeneratorHandler(w http.ResponseWriter, r *http.Request) {
	companyID := r.URL.Query().Get("id")
	title := r.URL.Query().Get("title")
	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")

	log.Println("dateFrom: ", dateFrom)
	log.Println("dateTo: ", dateTo)

	if companyID == "" || title == "" || dateFrom == "" || dateTo == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Генерация файла
	log.Println("FileGeneratorHandler was started")
	err := StartFileGenerator(companyID, title, dateFrom, dateTo)
	if err != nil {
		log.Println("StartFileGenerator failed: ", err)
		http.Error(w, "Ошибка при генерации файла: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Путь к сгенерированному файлу
	filePath := "/root/bitrixChatgpt/tables.docx"

	// Открытие файла для чтения
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Ошибка при открытии файла: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Настройка заголовков для скачивания файла
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")

	// Отправка файла на клиент
	http.ServeFile(w, r, filePath)
}

// StartFileGenerator запускает процесс генерации файла
func StartFileGenerator(companyId, title, dateFrom, dateTo string) error {
	var entries []DocumentEntry

	// Парсинг дат из строки (из формата "2006-01-02")
	startDate, err := time.Parse("2006-01-02", dateFrom)
	if err != nil {
		return fmt.Errorf("error parsing start date: %v", err)
	}

	endDate, err := time.Parse("2006-01-02", dateTo)
	if err != nil {
		return fmt.Errorf("error parsing end date: %v", err)
	}

	log.Println("StartFileGenerator started")
	// Форматирование дат
	formattedDateFrom := startDate.Format("02.01.2006")
	formattedDateTo := endDate.Format("02.01.2006")
	log.Println("Date completed")
	log.Println("formattedDateFrom: ", formattedDateFrom, "to: ", formattedDateTo)

	// Получаем данные об оплатах
	processesData, err := GetProcessesList(authorize.GlobalAuthId, companyId, "788")
	if err != nil {
		return fmt.Errorf("error getting processes list: %v", err)
	}
	AddProcessesToTable(&entries, filterProcessesByDate(processesData, startDate, endDate))

	// Получаем данные об актах
	itemsData, err := smart_processes.GetItemByCompany(companyId)
	if err != nil {
		return fmt.Errorf("error getting items by company: %v", err)
	}
	AddItemsToTable(&entries, filterItemsByDate(itemsData, startDate, endDate))

	// Сортировка записей по дате
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Date.Before(entries[j].Date)
	})

	// Формирование таблицы
	var tableDataLeft [][]string
	for _, entry := range entries {
		if entry.Type == "оплата" {
			office.AddToTableDataLeft(&tableDataLeft, "", entry.Description, "", entry.Credit)
		} else if entry.Type == "акт" {
			office.AddToTableDataLeft(&tableDataLeft, "", entry.Description, entry.Debit, "")
		}
	}

	var requisitesShortCompanyName string // сокращенное наименование компании
	requisitesList, err := requisites.GetRequisitesByCompanyID(companyId, authorize.GlobalAuthId)
	log.Println("requisitesList: ", requisitesList)
	if err != nil {
		log.Println("Error getting requisites by company:", err)
	}
	for _, req := range requisitesList {
		if req.RQCOMPANYNAME != "" {
			requisitesShortCompanyName = req.RQCOMPANYNAME
			break
		}
	}

	// Создаем документ Word
	err = office.StartWord(
		formattedDateFrom,          // dateFrom
		formattedDateTo,            // dateTo
		"",                         // dateForInvoice
		"",                         // dateForReturnPayment
		"",                         // numberOfCompletedWorks
		"",                         // dateOfCompletedWorks
		formattedDateFrom,          // dateForSaldo
		"",                         // sumOfSaldo
		requisitesShortCompanyName, // Используем переданный title в качестве secondCompany
		tableDataLeft,
	)
	if err != nil {
		return fmt.Errorf("error creating word document: %v", err)
	}
	return nil
}

// filterProcessesByDate фильтрует данные процессов по дате
func filterProcessesByDate(processes []ProcessesResponse, startDate, endDate time.Time) []ProcessesResponse {
	var filtered []ProcessesResponse
	for _, process := range processes {
		if process.Property638 != "" {
			processDate, err := time.Parse("02.01.2006", process.Property638)
			if err == nil && (processDate.Equal(startDate) || processDate.After(startDate)) && (processDate.Equal(endDate) || processDate.Before(endDate)) {
				filtered = append(filtered, process)
			}
		}
	}
	return filtered
}

// filterItemsByDate фильтрует данные элементов по дате
func filterItemsByDate(items []smart_processes.ItemsResponse, startDate, endDate time.Time) []smart_processes.ItemsResponse {
	var filtered []smart_processes.ItemsResponse
	for _, item := range items {
		itemDate, err := time.Parse("02.01.2006", item.UFCrm1712128088)
		if err == nil && (itemDate.Equal(startDate) || itemDate.After(startDate)) && (itemDate.Equal(endDate) || itemDate.Before(endDate)) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
