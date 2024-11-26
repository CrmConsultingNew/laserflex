package laserflex

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"strconv"
)

type Product struct {
	Name         string  // Наименование
	Quantity     float64 // Количество
	Price        float64 // Цена
	ImageBase64  string  // Вид детали (Base64)
	Material     float64 // Материал
	Laser        float64 // Лазер
	Bend         int     // Гиб
	Weld         int     // Свар
	Paint        int     // Окр
	Threading    int     // Резьба
	Countersink  int     // Зенк
	Drilling     int     // Сверл
	Rolling      int     // Вальц
	AddP         float64 // Допы П
	AddL         float64 // Допы Л
	PipeCutting  float64 // Труборез
	Construction float64 // Констр
	Delivery     float64 // Доставка
	Area         float64 // S, кв. м.
	Color        int     // Цвет
	P            string  // П
}

//H + (J-O) одна колонка - Раздел производство. (Суммируем в одно название - Производство)
//P это максимум , дальше не нужно

// ReadXlsProductRow читает только вторую строку из Excel файла
func ReadXlsProductRow(filename string) (*Product, error) {
	fmt.Println("xlsproduct started...")
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	// Читаем строку 2 из листа "Статистика"
	rows, err := f.GetRows("Статистика")
	if err != nil {
		return nil, fmt.Errorf("error reading rows: %v", err)
	}

	if len(rows) < 2 {
		return nil, errors.New("second row not found")
	}

	// Вторая строка (индекс 1)
	cells := rows[1]

	// Проверяем, что колонок достаточно
	if len(cells) < 21 {
		return nil, errors.New("not enough columns in the second row")
	}

	// Функция для преобразования строки в float64
	parseFloat := func(s string) float64 {
		v, _ := strconv.ParseFloat(s, 64)
		return v
	}

	// Функция для преобразования строки в int
	parseInt := func(s string) int {
		v, _ := strconv.Atoi(s)
		return v
	}

	// Сохраняем изображение из ячейки D2
	imageFile := "output_image.png"
	imageBase64 := ""
	if err := saveImageFromExcel(f, "Статистика", "D2"); err == nil {
		imgData, err := os.ReadFile(imageFile)
		if err == nil {
			imageBase64 = base64.StdEncoding.EncodeToString(imgData)
		}
	}

	// Создаём объект Product
	product := &Product{
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
		PipeCutting:  parseFloat(cells[15]),
		Construction: parseFloat(cells[16]),
		Delivery:     parseFloat(cells[17]),
		Area:         parseFloat(cells[18]),
		Color:        parseInt(cells[19]),
		P:            cells[20],
	}
	fmt.Println("xlsproduct finished...")
	return product, nil
}

func saveImageFromExcel(f *excelize.File, sheet, cell string) error {
	// Получаем массив картинок из ячейки
	pictures, err := f.GetPictures(sheet, cell)
	if err != nil {
		return fmt.Errorf("error extracting image from cell %s: %v", cell, err)
	}

	// Проверяем, есть ли хотя бы одна картинка
	if len(pictures) == 0 {
		return fmt.Errorf("no images found in cell %s", cell)
	}

	// Извлекаем данные первого изображения
	//imageData := pictures[0].File

	// Сохраняем данные изображения в файл
	/*if err := os.WriteFile(outputPath, imageData, 0644); err != nil {
		return fmt.Errorf("error saving image to file: %v", err)
	}*/
	return nil
}
