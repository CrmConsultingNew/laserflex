package laserflex

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
)

type Registry struct {
	CountOfMaterials []int    `json:"count_of_materials,omitempty"`  // Количество материала
	LaserWorks       []string `json:"laser_works,omitempty"`         // Лазерные работы
	PipeCutting      []string `json:"pipe_cutting,omitempty"`        // Труборез
	BendWorks        []int    `json:"bend_works,omitempty"`          // Гибочные работы
	TimeOfLaserWorks []int    `json:"time_of_laser_works,omitempty"` // Время лазерных работ
	Production       []string `json:"production,omitempty"`          // Производство
	Paint            []string `json:"paint,omitempty"`               // Нанесение покрытий
}

// Преобразует строку в `int`, игнорируя пустые значения
func parseStringToInt(input string) int {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0
	}

	value, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("Warning: unable to parse int from string '%s': %v\n", input, err)
		return 0
	}
	return value
}

// Функция для чтения реестра из Excel
func ReadXlsRegistryRows(filename string) (Registry, error) {
	fmt.Println("Processing Registry Excel file...")

	f, err := excelize.OpenFile(filename)
	if err != nil {
		return Registry{}, fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Реестр")
	if err != nil {
		return Registry{}, fmt.Errorf("error reading rows: %v", err)
	}

	// Определяем начальный и конечный столбцы
	startCol, endCol := -1, -1
	for i, cell := range rows[0] { // Считываем заголовки
		if strings.Contains(cell, "№ заказа") {
			startCol = i
		}
		if strings.Contains(cell, "Дата сдачи") {
			endCol = i
		}
	}

	if startCol == -1 || endCol == -1 {
		return Registry{}, fmt.Errorf("unable to find required columns (№ заказа or Дата сдачи)")
	}

	// Карты для хранения значений
	registry := Registry{
		CountOfMaterials: []int{},
		LaserWorks:       []string{},
		PipeCutting:      []string{},
		BendWorks:        []int{},
		TimeOfLaserWorks: []int{},
		Production:       []string{},
		Paint:            []string{},
	}

	// Регулярные выражения для сопоставления заголовков с полями структуры
	headerMapping := map[string]string{
		"Количество материала": "CountOfMaterials",
		"Лазерные работы":      "LaserWorks",
		"Труборез":             "PipeCutting",
		"Гибочные работы":      "BendWorks",
		"Время лазерных работ": "TimeOfLaserWorks",
		"Производство":         "Production",
		"Нанесение покрытий":   "Paint",
	}

	// Сопоставляем заголовки с полями структуры
	columnFieldMap := make(map[int]string)
	for i := startCol; i <= endCol; i++ {
		header := rows[0][i]
		for key, field := range headerMapping {
			if strings.Contains(header, key) {
				columnFieldMap[i] = field
				break
			}
		}
	}

	// Считываем строки таблицы
	for rowIndex, cells := range rows[1:] {
		isEmptyRow := true

		// Проверяем, что вся строка пустая в диапазоне startCol до endCol
		for colIndex := startCol; colIndex <= endCol && colIndex < len(cells); colIndex++ {
			if cells[colIndex] != "" {
				isEmptyRow = false
				break
			}
		}

		// Если строка полностью пустая, заканчиваем обработку
		if isEmptyRow {
			fmt.Printf("Empty row found at index %d, stopping processing.\n", rowIndex+1)
			break
		}

		// Заполняем значения в структуру
		for colIndex := startCol; colIndex <= endCol && colIndex < len(cells); colIndex++ {
			cellValue := cells[colIndex]
			fieldName, exists := columnFieldMap[colIndex]
			if !exists {
				continue
			}

			switch fieldName {
			case "CountOfMaterials":
				registry.CountOfMaterials = append(registry.CountOfMaterials, parseStringToInt(cellValue))
			case "LaserWorks":
				registry.LaserWorks = append(registry.LaserWorks, cellValue)
			case "PipeCutting":
				registry.PipeCutting = append(registry.PipeCutting, cellValue)
			case "BendWorks":
				registry.BendWorks = append(registry.BendWorks, parseStringToInt(cellValue))
			case "TimeOfLaserWorks":
				registry.TimeOfLaserWorks = append(registry.TimeOfLaserWorks, parseStringToInt(cellValue))
			case "Production":
				registry.Production = append(registry.Production, cellValue)
			case "Paint":
				registry.Paint = append(registry.Paint, cellValue)
			}
		}
	}

	// Выводим данные в консоль
	fmt.Printf("\nFinal Registry Data:\n")
	fmt.Printf("CountOfMaterials: %v\n", registry.CountOfMaterials)
	fmt.Printf("LaserWorks: %v\n", registry.LaserWorks)
	fmt.Printf("PipeCutting: %v\n", registry.PipeCutting)
	fmt.Printf("BendWorks: %v\n", registry.BendWorks)
	fmt.Printf("TimeOfLaserWorks: %v\n", registry.TimeOfLaserWorks)
	fmt.Printf("Production: %v\n", registry.Production)
	fmt.Printf("Paint: %v\n", registry.Paint)

	fmt.Println("Registry Excel processing completed.")
	return registry, nil
}
