package laserflex

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
)

func ParseFile() {
	filePath := "КП тест.xlsx"
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Printf("Ошибка открытия файла: %s\n", err)
		return
	}

	// Указываем лист и искомое значение
	sheetName := f.GetSheetName(1)
	searchValue := "Заголовок таблицы" // Здесь заменить на искомое значение

	// Находим ячейку с искомым значением
	startRow, startCol, err := findCell(f, sheetName, searchValue)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Найдено значение '%s' в ячейке [%d, %d]\n", searchValue, startRow+1, startCol+1)

	// Определяем границы подтаблицы
	rangeStr, err := detectSubtable(f, sheetName, startRow, startCol)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Диапазон подтаблицы: %s\n", rangeStr)
}

func findCell(f *excelize.File, sheetName string, searchValue string) (int, int, error) {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return -1, -1, err
	}

	for r, row := range rows {
		for c, cell := range row {
			if strings.TrimSpace(cell) == searchValue {
				return r, c, nil
			}
		}
	}

	return -1, -1, fmt.Errorf("значение '%s' не найдено", searchValue)
}

// Определение границ подтаблицы
func detectSubtable(f *excelize.File, sheetName string, startRow, startCol int) (string, error) {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return "", err
	}

	// Определяем конец по горизонтали (строка)
	endCol := startCol
	for col := startCol; col < len(rows[startRow]); col++ {
		if strings.TrimSpace(rows[startRow][col]) == "" {
			break
		}
		endCol = col
	}

	// Определяем конец по вертикали (столбец)
	endRow := startRow
	for row := startRow; row < len(rows); row++ {
		empty := true
		for col := startCol; col <= endCol; col++ {
			if strings.TrimSpace(rows[row][col]) != "" {
				empty = false
				break
			}
		}
		if empty {
			break
		}
		endRow = row
	}

	// Формируем диапазон подтаблицы
	startCell, _ := excelize.CoordinatesToCellName(startCol+1, startRow+1)
	endCell, _ := excelize.CoordinatesToCellName(endCol+1, endRow+1)
	return fmt.Sprintf("%s:%s", startCell, endCell), nil
}
