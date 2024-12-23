package laserflex

import (
	"github.com/xuri/excelize/v2"
	"log"
	"strings"
)

func CheckCoatingColumn(fileName string) bool {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		log.Printf("Error opening Excel file: %v", err)
		return false
	}
	defer f.Close()

	rows, err := f.GetRows("Реестр")
	if err != nil {
		log.Printf("Error reading Excel rows: %v", err)
		return false
	}

	// Определяем индекс столбца "Нанесение покрытий"
	var coatingColumnIndex int = -1
	for i, cell := range rows[0] {
		if cell == "Нанесение покрытий" {
			coatingColumnIndex = i
			break
		}
	}

	if coatingColumnIndex == -1 {
		log.Println("Column 'Нанесение покрытий' not found")
		return false
	}

	// Проверяем наличие заполненных ячеек и выводим их
	var foundValues []string
	for _, row := range rows[1:] {
		if len(row) > coatingColumnIndex {
			value := strings.TrimSpace(row[coatingColumnIndex])
			if value != "" {
				foundValues = append(foundValues, value)
			}
		}
	}

	if len(foundValues) > 0 {
		log.Printf("Found values in 'Нанесение покрытий': %v", foundValues)
		return true
	}

	log.Println("No values found in 'Нанесение покрытий'")
	return false
}
