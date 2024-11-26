package laserflex

import (
	"github.com/xuri/excelize/v2"
	"log"
)

func ReadXlsProducts(filename string) map[string][]string {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		log.Println("Error opening file:", err)
		return nil
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("Error closing file:", err)
		}
	}()

	// Для хранения данных, где ключ — значение из первой ячейки строки
	data := make(map[string][]string)

	// Получаем все строки из листа "КП"
	rows, err := f.GetRows("КП")
	if err != nil {
		log.Println("Error getting rows:", err)
		return nil
	}

	// Перебираем строки
	for rowIndex, row := range rows {
		// Пропускаем первую строку (например, заголовки)
		if rowIndex == 0 {
			continue
		}

		// Если строка пустая, пропускаем её
		if len(row) == 0 {
			continue
		}

		// Первый элемент строки становится ключом
		key := row[0]
		// Остальные элементы добавляются в значение
		if len(row) > 1 {
			data[key] = row[1:]
		} else {
			data[key] = []string{}
		}
	}

	return data
}
