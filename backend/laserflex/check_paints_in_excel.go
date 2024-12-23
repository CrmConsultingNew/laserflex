package laserflex

import (
	"github.com/xuri/excelize/v2"
	"log"
	"strings"
)

func checkCoatingColumn(fileName string) bool {
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

	// Проверяем наличие заполненных ячеек
	for rowIndex, row := range rows[1:] {
		if len(row) > coatingColumnIndex && strings.TrimSpace(row[coatingColumnIndex]) != "" {
			log.Printf("Found value in 'Нанесение покрытий' at row %d: %s", rowIndex+2, row[coatingColumnIndex])
			return true
		}
	}
	log.Println("No filled cells found in column 'Нанесение покрытий'")
	return false
}
