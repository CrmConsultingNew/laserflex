package main

import (
	"bitrix_app/backend/bitrix/endpoints"
	"bitrix_app/backend/routes"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func ReadXlsProducts(filename string) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Get value from cell by given worksheet name and cell reference.
	cell, err := f.GetCellValue("Статистика", "A2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cell)
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("КП")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}

func main() {

	//ReadXlsProducts("file_downloaded_xls1.xlsx")
	//return

	fmt.Println("service starting...")

	// Загрузка переменных окружения из файла .env
	if err := godotenv.Load(filepath.Join(".env")); err != nil {
		log.Print("No .env file found")
	} else {
		fmt.Println("Loaded .env file")
	}

	// Установка домена для Bitrix24
	endpoints.BitrixDomain = os.Getenv("BITRIX_DOMAIN")
	endpoints.NewBitrixDomain = os.Getenv("NEW_BITRIX_DOMAIN")

	// Инициализация маршрутов
	routes.Router()

	// Запуск сервера
	server := &http.Server{
		Addr:              ":9090",
		ReadHeaderTimeout: 3 * time.Second,
	}

	fmt.Printf("server started on addr: %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Server started with error")
		panic(err)
	}
}
