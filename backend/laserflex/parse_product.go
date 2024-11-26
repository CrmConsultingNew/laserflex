package laserflex

import (
	"encoding/base64"
	"fmt"
	"github.com/xuri/excelize/v2"
	"regexp"
	"strconv"
	"strings"
)

type Product struct {
	Name        string  // Наименование (столбец A)
	Quantity    float64 // Количество (столбец B)
	Price       float64 // Цена (столбец C)
	ImageBase64 string  // Изображение в Base64 (столбец D)
	Material    float64 // Материал (столбец E)
	Laser       float64 // Лазер (столбец F)
	Bend        int     // Гиб (столбец G)
	Weld        int     // Свар (столбец H)
	Paint       int     // Окраска (столбец I)
	Production  float64 // Производство (Сумма столбцов H, J, K, L, M, N, O)
	AddP        float64 // Допы П (столбец N)
	AddL        float64 // Допы Л (столбец O)
	PipeCutting float64 // Труборез (столбец P)
}

// Функция для конверсии строки с запятой в float64
func parseLocalizedFloat(s string) float64 {
	// Удаляем пробелы только между цифрами
	re := regexp.MustCompile(`(\d)\s+(\d)`)
	s = re.ReplaceAllString(s, "$1$2")

	// Заменяем запятую на точку
	s = strings.ReplaceAll(s, ",", ".")

	// Парсим строку в float64
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Printf("Warning: unable to parse float from string '%s': %v\n", s, err)
		return 0 // Возвращаем 0, если конвертация не удалась
	}
	return v
}

// Функция для чтения строк из Excel
func ReadXlsProductRows(filename string) ([]Product, error) {
	fmt.Println("Processing Excel file...")

	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Статистика")
	if err != nil {
		return nil, fmt.Errorf("error reading rows: %v", err)
	}

	var products []Product

	// Преобразование строки в int
	parseInt := func(s string) int {
		v, _ := strconv.Atoi(s)
		return v
	}

	// Обработка каждой строки
	for i, cells := range rows {
		if i == 0 || len(cells) < 16 { // Пропускаем заголовок и проверяем минимальное количество столбцов
			continue
		}

		// Завершаем обработку, если первая колонка пустая
		if cells[0] == "" {
			break
		}

		// Логируем данные строки
		fmt.Printf("Row %d: %+v\n", i+1, cells)

		// Получение Base64 строки изображения из ячейки
		imageBase64 := ""
		imageData, err := getImageBase64FromExcel(f, "Статистика", fmt.Sprintf("D%d", i+1))
		if err == nil {
			imageBase64 = imageData
		} else {
			fmt.Printf("Warning: unable to get image for row %d: %v\n", i+1, err)
		}

		// Суммируем для Production
		production := parseLocalizedFloat(cells[7]) + // H
			parseLocalizedFloat(cells[9]) + // J
			parseLocalizedFloat(cells[10]) + // K
			parseLocalizedFloat(cells[11]) + // L
			parseLocalizedFloat(cells[12]) + // M
			parseLocalizedFloat(cells[13]) + // N
			parseLocalizedFloat(cells[14]) // O

		product := Product{
			Name:        cells[0],
			Quantity:    parseLocalizedFloat(cells[1]),
			Price:       parseLocalizedFloat(cells[2]),
			ImageBase64: imageBase64,
			Material:    parseLocalizedFloat(cells[4]),
			Laser:       parseLocalizedFloat(cells[5]),
			Bend:        parseInt(cells[6]),
			Weld:        parseInt(cells[7]),
			Paint:       parseInt(cells[8]),
			Production:  production,
			AddP:        parseLocalizedFloat(cells[13]),
			AddL:        parseLocalizedFloat(cells[14]),
			PipeCutting: parseLocalizedFloat(cells[15]),
		}

		fmt.Printf("Parsed Product: %+v\n", product)

		products = append(products, product)
	}

	fmt.Println("Excel processing completed.")
	return products, nil
}

// getImageBase64FromExcel извлекает изображение из ячейки и возвращает его в виде строки Base64
func getImageBase64FromExcel(f *excelize.File, sheet, cell string) (string, error) {
	pictures, err := f.GetPictures(sheet, cell)
	if err != nil {
		return "", fmt.Errorf("error extracting image from cell %s: %v", cell, err)
	}

	if len(pictures) == 0 {
		return "", fmt.Errorf("no images found in cell %s", cell)
	}

	// Извлекаем данные первого изображения
	imageData := pictures[0].File

	// Кодируем изображение в Base64
	return base64.StdEncoding.EncodeToString(imageData), nil
}
