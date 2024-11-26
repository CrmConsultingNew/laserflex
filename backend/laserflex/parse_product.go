package laserflex

import (
	"encoding/base64"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"strconv"
)

type Product struct {
	Name         string  // Наименование (столбец A)
	Quantity     float64 // Количество (столбец B)
	Price        float64 // Цена (столбец C)
	ImageBase64  string  // Изображение в Base64 (столбец D)
	Material     float64 // Материал (столбец E)
	Laser        float64 // Лазер (столбец F)
	Bend         int     // Гиб (столбец G)
	Weld         int     // Свар (столбец H)
	Paint        int     // Окраска (столбец I)
	Threading    int     // Резьба (столбец J)
	Countersink  int     // Зенк (столбец K)
	Drilling     int     // Сверление (столбец L)
	Rolling      int     // Вальцовка (столбец M)
	AddP         float64 // Допы П (столбец N)
	AddL         float64 // Допы Л (столбец O)
	Production   float64 // Производство (Сумма столбцов G-N)
	PipeCutting  float64 // Труборез (столбец P)
	Construction float64 // Конструирование (столбец Q)
	Delivery     float64 // Доставка (столбец R)
	Area         float64 // Площадь (S, столбец S)
	Color        int     // Цвет (столбец T)
	P            string  // P (столбец U)
}

//H + (J-O) одна колонка - Раздел производство. (Суммируем в одно название - Производство)
//P это максимум , дальше не нужно

// ReadXlsProductRow читает только вторую строку из Excel файла
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
		if i == 0 || len(cells) < 21 { // Пропускаем заголовок и проверяем минимальное количество столбцов
			continue
		}

		// Завершаем обработку, если первая колонка пустая
		if cells[0] == "" {
			break
		}

		// Сохраняем изображение
		imageFile := fmt.Sprintf("output_image_row_%d.png", i)
		imageBase64 := ""
		if err := saveImageFromExcel(f, "Статистика", fmt.Sprintf("D%d", i+1)); err == nil {
			imgData, err := os.ReadFile(imageFile)
			if err == nil {
				imageBase64 = base64.StdEncoding.EncodeToString(imgData)
			}
		}

		// Создаём объект Product
		product := Product{
			Name:         cells[0],
			Quantity:     parseFloat(cells[1]),
			Price:        parseFloat(cells[2]),
			ImageBase64:  imageBase64,
			Material:     parseFloat(cells[4]),
			Laser:        parseFloat(cells[5]),
			Bend:         parseInt(cells[6]),
			Weld:         parseInt(cells[7]),
			Paint:        parseInt(cells[8]),
			Threading:    parseInt(cells[9]),
			Countersink:  parseInt(cells[10]),
			Drilling:     parseInt(cells[11]),
			Rolling:      parseInt(cells[12]),
			AddP:         parseFloat(cells[13]),
			AddL:         parseFloat(cells[14]),
			Production:   parseFloat(cells[6]) + parseFloat(cells[7]) + parseFloat(cells[8]) + parseFloat(cells[9]) + parseFloat(cells[10]) + parseFloat(cells[11]) + parseFloat(cells[12]),
			PipeCutting:  parseFloat(cells[15]),
			Construction: parseFloat(cells[16]),
			Delivery:     parseFloat(cells[17]),
			Area:         parseFloat(cells[18]),
			Color:        parseInt(cells[19]),
			P:            cells[20],
		}

		products = append(products, product)
	}

	fmt.Println("Excel processing completed.")
	return products, nil
}

func saveImageFromExcel(f *excelize.File, sheet, cell string) error {
	// Получение изображений из ячейки
	pictures, err := f.GetPictures(sheet, cell)
	if err != nil {
		return fmt.Errorf("error extracting image from cell %s: %v", cell, err)
	}

	if len(pictures) == 0 {
		return fmt.Errorf("no images found in cell %s", cell)
	}

	// Пример сохранения первого изображения
	imageData := pictures[0].File
	outputPath := "output_image.png"
	if err := os.WriteFile(outputPath, imageData, 0644); err != nil {
		return fmt.Errorf("error saving image to file: %v", err)
	}

	return nil
}
