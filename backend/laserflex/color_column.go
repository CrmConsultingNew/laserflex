package laserflex

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"strings"
)

func parseSheetForColorColumnAndTasks(fileName string) (string, string, []string, error) {
	// Открываем файл
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return "", "", nil, fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	// Переходим на нужный лист
	sheetName := "Реестр"
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return "", "", nil, fmt.Errorf("error reading rows from sheet '%s': %v", sheetName, err)
	}

	// Нормализуем заголовки для поиска
	colorHeader := strings.ToLower("Цвет/цинк")
	orderHeader := "№ заказа"
	customerHeader := "Заказчик"

	var colorColumnIndex, headerRowIndex int = -1, -1
	var headers = map[string]int{
		orderHeader:    -1,
		customerHeader: -1,
	}

	// Ищем заголовки и определяем индексы
	for i, row := range rows {
		for j, cell := range row {
			normalizedCell := strings.TrimSpace(strings.ToLower(cell))
			if normalizedCell == colorHeader {
				colorColumnIndex = j
				headerRowIndex = i
				log.Printf("Header '%s' found in row %d, column %d\n", colorHeader, i+1, j+1)
			}
			if _, ok := headers[cell]; ok {
				headers[cell] = j
			}
		}
		if colorColumnIndex != -1 && headers[orderHeader] != -1 && headers[customerHeader] != -1 {
			break
		}
	}

	// Проверяем наличие всех необходимых заголовков
	if colorColumnIndex == -1 {
		return "", "", nil, fmt.Errorf("column '%s' not found", colorHeader)
	}
	for header, index := range headers {
		if index == -1 {
			return "", "", nil, fmt.Errorf("missing required header: %s", header)
		}
	}

	// Собираем значения из найденного столбца, начиная с строки после заголовка
	var colorValues []string
	var orderNumber, customer string

	for i := headerRowIndex + 1; i < len(rows); i++ {
		row := rows[i]

		// Извлекаем первое значение из столбцов "№ Заказа" и "Заказчик"
		if orderNumber == "" && len(row) > headers[orderHeader] {
			orderNumber = strings.TrimSpace(row[headers[orderHeader]])
		}
		if customer == "" && len(row) > headers[customerHeader] {
			customer = strings.TrimSpace(row[headers[customerHeader]])
		}

		// Проверяем, чтобы индекс столбца не превышал длину строки
		if len(row) > colorColumnIndex {
			colorValue := strings.TrimSpace(row[colorColumnIndex])
			if colorValue != "" {
				colorValues = append(colorValues, colorValue)
			}
		}
	}

	// Проверяем, чтобы значения "№ Заказа" и "Заказчик" были найдены
	if orderNumber == "" || customer == "" {
		return "", "", nil, fmt.Errorf("missing required data: order number or customer")
	}

	return orderNumber, customer, colorValues, nil
}
