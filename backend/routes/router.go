package routes

import (
	"bitrix_app/backend/laserflex"
	"bitrix_app/backend/laserflex/authorize_backup"
	"net/http"
)

func Router() {

	//http.HandleFunc("/laser_checklist", laserflex.HandlerProcessProducts)

	http.HandleFunc("/laser_auth", authorize_backup.AuthorizeEndpoint)

	http.HandleFunc("/send_file", laserflex.LaserflexGetFile)
	http.HandleFunc("/files", laserflex.ListFilesHandler)        // Страница со списком файлов
	http.HandleFunc("/download/", laserflex.DownloadFileHandler) // Скачивание файла

}
