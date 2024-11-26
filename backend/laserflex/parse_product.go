package laserflex

import (
	"encoding/base64"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
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

// ReadXlsProductRows читает строки из Excel файла и возвращает массив структур Product
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

	// Преобразование строки в float64
	parseFloat := func(s string) float64 {
		v, _ := strconv.ParseFloat(s, 64)
		return v
	}

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

		// Получение Base64 строки изображения из ячейки
		imageBase64 := ""
		imageData, err := getImageBase64FromExcel(f, "Статистика", fmt.Sprintf("D%d", i+1))
		if err == nil {
			imageBase64 = imageData
		} else {
			fmt.Printf("Warning: unable to get image for row %d: %v\n", i+1, err)
		}

		// Суммируем для Production
		production := parseFloat(cells[7]) + // H
			parseFloat(cells[9]) + // J
			parseFloat(cells[10]) + // K
			parseFloat(cells[11]) + // L
			parseFloat(cells[12]) + // M
			parseFloat(cells[13]) + // N
			parseFloat(cells[14]) // O

		// Создаём объект Product
		product := Product{
			Name:        cells[0],
			Quantity:    parseFloat(cells[1]),
			Price:       parseFloat(cells[2]),
			ImageBase64: imageBase64,
			Material:    parseFloat(cells[4]),
			Laser:       parseFloat(cells[5]),
			Bend:        parseInt(cells[6]),
			Weld:        parseInt(cells[7]),
			Paint:       parseInt(cells[8]),
			Production:  production,
			AddP:        parseFloat(cells[13]),
			AddL:        parseFloat(cells[14]),
			PipeCutting: parseFloat(cells[15]),
		}

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
