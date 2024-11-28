package laserflex

import (
	"encoding/base64"
	"fmt"
	"github.com/xuri/excelize/v2"
	"math"
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
	Bend        float64 // Гиб (столбец G)
	Weld        float64 // Свар (столбец H)
	Paint       float64 // Окраска (столбец I)
	Production  float64 // Производство (Сумма столбцов H, J, K, L, M, N, O)
	AddP        float64 // Допы П (столбец N)
	AddL        float64 // Допы Л (столбец O)
	PipeCutting float64 // Труборез (столбец P)
}

// Специальная функция для обработки столбца Price
func parsePrice(input string) float64 {
	fmt.Printf("Original Price input: '%s'\n", input)

	// Удаление пробелов между цифрами
	re := regexp.MustCompile(`(\d)\s+(\d)`)
	input = re.ReplaceAllString(input, "$1$2")
	fmt.Printf("After removing spaces: '%s'\n", input)

	// Замена всех запятых на точки
	input = strings.ReplaceAll(input, ",", ".")

	// Проверка на наличие более одной точки
	if strings.Count(input, ".") > 1 {
		// Оставляем только последнюю точку
		parts := strings.Split(input, ".")
		input = strings.Join(parts[:len(parts)-1], "") + "." + parts[len(parts)-1]
		fmt.Printf("After fixing dots: '%s'\n", input)
	}

	// Преобразование в float64
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Printf("Error parsing Price: %v\n", err)
		return 0
	}

	// Округляем до двух знаков после запятой
	value = math.Round(value*100) / 100
	fmt.Printf("Parsed and rounded Price: %f\n", value)

	return value
}

// Функция для обработки чисел в других столбцах
func parseFloatOrInt(input string) float64 {
	fmt.Printf("Original input: '%s'\n", input)

	// Убираем пробелы между цифрами
	re := regexp.MustCompile(`(\d)\s+(\d)`)
	input = re.ReplaceAllString(input, "$1$2")
	//fmt.Printf("After removing spaces: '%s'\n", input)

	// Заменяем запятую на точку
	input = strings.ReplaceAll(input, ",", ".")
	//fmt.Printf("After replacing commas: '%s'\n", input)

	// Пробуем преобразовать в float64
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Printf("Warning: unable to parse float or int from string '%s': %v\n", input, err)
		return 0
	}

	//fmt.Printf("Parsed float: %f\n", value)
	return value
}

// Функция для чтения строк из Excel
func ReadXlsProductRows(filename string) ([]Product, error) {
	fmt.Println("Processing Excel file...")

	// Открываем файл Excel
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	// Читаем строки из листа "Статистика"
	rows, err := f.GetRows("Статистика")
	if err != nil {
		return nil, fmt.Errorf("error reading rows: %v", err)
	}

	var products []Product

	// Обрабатываем каждую строку
	for i, cells := range rows {
		if i == 0 || len(cells) < 16 { // Пропускаем заголовок и проверяем минимальное количество столбцов
			continue
		}

		// Проверяем первую ячейку на условия завершения
		if len(cells) > 0 {
			name := strings.TrimSpace(cells[0]) // Убираем лишние пробелы
			if name == "" || strings.EqualFold(name, "Доставка") || strings.Contains(strings.ToLower(name), "общее") {
				fmt.Printf("Terminating parsing at row %d: Name='%s'\n", i+1, name)
				break
			}
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
		production := 0.0
		for _, colIndex := range []int{7, 9, 10, 11, 12, 13, 14} {
			if colIndex < len(cells) && cells[colIndex] != "" {
				production += parseFloatOrInt(cells[colIndex])
			} else {
				fmt.Printf("Warning: missing or empty value at column %d, row %d\n", colIndex, i+1)
			}
		}

		// Создаём объект Product
		product := Product{
			Name:        cells[0],
			Quantity:    parseFloatOrInt(cells[1]),
			Price:       parsePrice(cells[2]), // Используем parsePrice для обработки столбца C
			ImageBase64: imageBase64,
			Material:    parseFloatOrInt(cells[4]),
			Laser:       parseFloatOrInt(cells[5]),
			Bend:        parseFloatOrInt(cells[6]),
			Weld:        parseFloatOrInt(cells[7]),
			Paint:       parseFloatOrInt(cells[8]),
			Production:  production,
			AddP:        parseFloatOrInt(cells[13]),
			AddL:        parseFloatOrInt(cells[14]),
			PipeCutting: parseFloatOrInt(cells[15]),
		}

		products = append(products, product)
		fmt.Printf("Parsed Product: %+v\n", product)
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
