package laserflex

import (
	"github.com/xuri/excelize/v2"
	"log"
	"strings"
)

func ParseSheetForColorColumn(fileName string) []string {
	// Открываем файл
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()

	// Переходим на нужный лист
	sheetName := "Реестр"
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatalf("Error reading rows from sheet '%s': %v", sheetName, err)
	}

	// Нормализуем заголовки для поиска
	targetHeader := strings.ToLower("Цвет/цинк")

	var colorColumnIndex int = -1
	var headerRowIndex int = -1

	// Ищем заголовок и определяем индекс
	for i, row := range rows {
		for j, cell := range row {
			if strings.TrimSpace(strings.ToLower(cell)) == targetHeader {
				colorColumnIndex = j
				headerRowIndex = i
				log.Printf("Header '%s' found in row %d, column %d\n", targetHeader, i+1, j+1)
				break
			}
		}
		if colorColumnIndex != -1 {
			break
		}
	}

	// Если не найдено, выводим сообщение и завершение
	if colorColumnIndex == -1 {
		log.Println("Column 'Цвет/цинк' not found")
		return nil
	}

	// Собираем значения из найденного столбца, начиная с строки после заголовка
	var colorValues []string
	for i := headerRowIndex + 1; i < len(rows); i++ {
		row := rows[i]
		// Проверяем, чтобы индекс столбца не превышал длину строки
		if len(row) > colorColumnIndex {
			value := strings.TrimSpace(row[colorColumnIndex])
			if value != "" {
				colorValues = append(colorValues, value)
			}
		}
	}

	// Выводим значения в консоль
	log.Printf("Values from column 'Цвет/цинк' below the header: %v\n", colorValues)
	return colorValues
}
