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

	var products []Product
	const sheet = "Статистика"

	// Определяем количество строк
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("error reading rows: %v", err)
	}

	numRows := len(rows)

	// Обрабатываем каждую строку
	for i := 2; i <= numRows; i++ { // Начинаем с 2, т.к. 1 - заголовок
		name, _ := f.GetCellValue(sheet, fmt.Sprintf("A%d", i))
		if strings.TrimSpace(name) == "" {
			fmt.Printf("Skipping empty row at %d\n", i)
			continue
		}
		if strings.Contains(strings.ToLower(name), "общее") {
			fmt.Printf("Terminating parsing at row %d: Name='%s'\n", i, name)
			break
		}

		imageBase64 := ""
		if imgData, err := getImageBase64FromExcel(f, sheet, fmt.Sprintf("D%d", i)); err == nil {
			imageBase64 = imgData
		} else {
			fmt.Printf("Warning: unable to get image for row %d: %v\n", i, err)
		}

		production := 0.0
		for _, col := range []string{"H", "J", "K", "L", "M", "N", "O"} {
			val, _ := f.GetCellValue(sheet, fmt.Sprintf("%s%d", col, i))
			production += parseFloatOrInt(strings.TrimSpace(val))
		}

		// Получаем и преобразуем значения столбцов
		quantityStr, _ := f.GetCellValue(sheet, fmt.Sprintf("B%d", i))
		priceStr, _ := f.GetCellValue(sheet, fmt.Sprintf("C%d", i))
		materialStr, _ := f.GetCellValue(sheet, fmt.Sprintf("E%d", i))
		laserStr, _ := f.GetCellValue(sheet, fmt.Sprintf("F%d", i))
		bendStr, _ := f.GetCellValue(sheet, fmt.Sprintf("G%d", i))
		weldStr, _ := f.GetCellValue(sheet, fmt.Sprintf("H%d", i))
		paintStr, _ := f.GetCellValue(sheet, fmt.Sprintf("I%d", i))
		addPStr, _ := f.GetCellValue(sheet, fmt.Sprintf("N%d", i))
		addLStr, _ := f.GetCellValue(sheet, fmt.Sprintf("O%d", i))
		pipeCuttingStr, _ := f.GetCellValue(sheet, fmt.Sprintf("P%d", i))

		product := Product{
			Name:        name,
			Quantity:    parseFloatOrInt(strings.TrimSpace(quantityStr)),
			Price:       parsePrice(strings.TrimSpace(priceStr)),
			ImageBase64: imageBase64,
			Material:    parseFloatOrInt(strings.TrimSpace(materialStr)),
			Laser:       parseFloatOrInt(strings.TrimSpace(laserStr)),
			Bend:        parseFloatOrInt(strings.TrimSpace(bendStr)),
			Weld:        parseFloatOrInt(strings.TrimSpace(weldStr)),
			Paint:       parseFloatOrInt(strings.TrimSpace(paintStr)),
			Production:  production,
			AddP:        parseFloatOrInt(strings.TrimSpace(addPStr)),
			AddL:        parseFloatOrInt(strings.TrimSpace(addLStr)),
			PipeCutting: parseFloatOrInt(strings.TrimSpace(pipeCuttingStr)),
		}

		products = append(products, product)
		fmt.Printf("Parsed Product: %+v\n", product)
	}

	fmt.Println("\nExcel processing completed. Parsed Products:")
	for _, product := range products {
		fmt.Printf("%+v\n", product)
	}

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

// Извлекает данные для указанных столбцов
func ExtractData(cells []string, headers map[string]int, columns []string) map[string]string {
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
func GetValue(cells []string, index int) string {
	if index >= 0 && index < len(cells) {
		return cells[index]
	}
	return ""
}
