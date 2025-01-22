package main

import (
	"bitrix_app/backend/bitrix/endpoints"
	"bitrix_app/backend/routes"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {

	//laserflex.StartJsonConverterFromExcel("file.xlsx")

	//laserflex.ReadXlsRegistryWithConditions("file.xlsx")

	//laserflex.ReadXlsProducts("file_downloaded_xls1.xlsx")

	/*	laserflex.ReadXlsProductRows("file_downloaded_xls1.xlsx")

		return*/

	/*fileName := "file.xlsx"
	if laserflex.CheckCoatingColumn(fileName) {
		// Если есть данные, получаем цвета из "Цвет/цинк"
		colors := laserflex.ParseSheetForColorColumn(fileName)
		_, err := laserflex.AddTaskToGroupColor("Проверить наличие ЛКП на складе в ОМТС", 149, 12, 1046, 791, colors)
		if err != nil {
			log.Printf("Error creating task with colors: %v", err)
			return
		}

	} else {
		// Если данных нет, создаём задачу без цветов
		_, err := laserflex.AddTaskToGroupColor("Задача в ОМТС с материалами из накладной", 149, 12, 1046, 791, nil)
		if err != nil {
			log.Printf("Error creating task without colors: %v", err)
			return
		}
	}

	return*/
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
