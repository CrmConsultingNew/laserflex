package main

import (
	"bitrix_app/backend/bitrix/endpoints"
	"bitrix_app/backend/laserflex"
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

	laserflex.ReadXlsRegistryRows("file.xlsx")

	/*laserflex.ReadXlsProducts("file_downloaded_xls1.xlsx")
	laserflex.ReadXlsProductRow("file_downloaded_xls1.xlsx")
	return*/

	return
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
