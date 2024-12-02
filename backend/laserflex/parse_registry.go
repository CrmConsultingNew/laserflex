package laserflex

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"strings"
)

// Универсальная структура для хранения данных из строки
type ParsedData struct {
	LaserWorks  *WorkGroup `json:"laser_works,omitempty"`  // Для строк с Лазерными работами
	PipeCutting *WorkGroup `json:"pipe_cutting,omitempty"` // Для строк с Труборезом
	BendWorks   *WorkGroup `json:"bend_works,omitempty"`   // Для строк с Гибочными работами
	Production  *WorkGroup `json:"production,omitempty"`   // Для строк с Производством
}

// WorkGroup структура для хранения данных группы
type WorkGroup struct {
	GroupID int               `json:"group_id"`
	Data    map[string]string `json:"data"`
}

// Функция для чтения таблицы и разделения данных по условиям
func ReadXlsRegistryWithConditions(filename string) ([]ParsedData, error) {
	fmt.Println("Processing Registry Excel file...")

	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Реестр")
	if err != nil {
		return nil, fmt.Errorf("error reading rows: %v", err)
	}

	// Определяем индексы заголовков
	headers := map[string]int{
		"№ заказа":             -1,
		"Заказчик":             -1,
		"Менеджер":             -1,
		"Количество материала": -1,
		"Лазерные работы":      -1,
		"Труборез":             -1,
		"Гибочные работы":      -1,
		"Время лазерных работ": -1,
		"Производство":         -1,
		"Нанесение покрытий":   -1, // Добавляем столбец "Нанесение покрытий"
		"Комментарий":          -1,
	}

	// Найдем индексы всех необходимых заголовков
	for i, cell := range rows[0] {
		for header := range headers {
			if strings.Contains(cell, header) {
				headers[header] = i
				break
			}
		}
	}

	// Проверим, что все необходимые заголовки найдены
	for header, index := range headers {
		if index == -1 {
			return nil, fmt.Errorf("missing required header: %s", header)
		}
	}

	// Парсинг строк
	var parsedRows []ParsedData
	for _, cells := range rows[1:] {
		// Проверяем, если строка полностью пустая
		isEmptyRow := true
		for _, index := range headers {
			if index < len(cells) && cells[index] != "" {
				isEmptyRow = false
				break
			}
		}
		if isEmptyRow {
			break
		}

		// Создаем структуру для хранения данных
		parsedRow := ParsedData{}

		// Проверяем условия для каждого столбца
		if value := getValue(cells, headers["Лазерные работы"]); value != "" {
			parsedRow.LaserWorks = &WorkGroup{
				GroupID: 1,
				Data: extractData(cells, headers, []string{
					"№ заказа", "Заказчик", "Менеджер", "Количество материала", "Лазерные работы", "Нанесение покрытий", "Комментарий",
				}),
			}
		}

		if value := getValue(cells, headers["Труборез"]); value != "" {
			parsedRow.PipeCutting = &WorkGroup{
				GroupID: 11,
				Data: extractData(cells, headers, []string{
					"№ заказа", "Заказчик", "Менеджер", "Количество материала", "Труборез", "Нанесение покрытий", "Комментарий",
				}),
			}
		}

		if value := getValue(cells, headers["Гибочные работы"]); value != "" {
			parsedRow.BendWorks = &WorkGroup{
				GroupID: 10,
				Data: extractData(cells, headers, []string{
					"№ заказа", "Заказчик", "Менеджер", "Количество материала", "Гибочные работы", "Нанесение покрытий", "Комментарий",
				}),
			}
		}

		if value := getValue(cells, headers["Производство"]); value != "" {
			parsedRow.Production = &WorkGroup{
				GroupID: 2,
				Data: extractData(cells, headers, []string{
					"№ заказа", "Заказчик", "Менеджер", "Количество материала", "Производство", "Нанесение покрытий", "Комментарий",
				}),
			}
		}

		// Если данные для текущей строки соответствуют хотя бы одному условию, добавляем строку
		if parsedRow.LaserWorks != nil || parsedRow.PipeCutting != nil || parsedRow.BendWorks != nil || parsedRow.Production != nil {
			parsedRows = append(parsedRows, parsedRow)
		}
	}

	// Выводим результат
	//fmt.Printf("\nParsed Rows:\n%v\n", parsedRows)

	return parsedRows, nil
}

func StartJsonConverterFromExcel(filename string) {
	parsedRows, err := ReadXlsRegistryWithConditions(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Пишем результат в JSON файл
	err = WriteParsedDataToJSON(parsedRows, "output.json")
	if err != nil {
		fmt.Println("Error writing to JSON file:", err)
		return
	}
}

// Извлекает данные для указанных столбцов
func extractData(cells []string, headers map[string]int, columns []string) map[string]string {
	data := make(map[string]string)
	for _, column := range columns {
		index := headers[column]
		if index >= 0 && index < len(cells) {
			data[column] = cells[index]
		}
	}
	return data
}

// Получает значение из ячейки или возвращает пустую строку, если индекс выходит за пределы
func getValue(cells []string, index int) string {
	if index >= 0 && index < len(cells) {
		return cells[index]
	}
	return ""
}

func WriteParsedDataToJSON(parsedRows []ParsedData, outputFile string) error {
	// Преобразуем данные в формат JSON с отступами для читабельности
	jsonData, err := json.MarshalIndent(parsedRows, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling data to JSON: %v", err)
	}

	// Создаем или перезаписываем файл
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Записываем JSON данные в файл
	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error writing JSON to file: %v", err)
	}

	fmt.Printf("Data successfully written to %s\n", outputFile)
	return nil
}
