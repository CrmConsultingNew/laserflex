package spreadsheets

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
	"os"
)

func GoogleSheetsUpdate(row, column int, value string) {
	data, err := os.ReadFile("local-abbey-420416-ecb7d8e898c0.json")
	CheckError(err)
	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	CheckError(err)
	client := conf.Client(context.TODO())

	service := spreadsheet.NewServiceWithClient(client)

	//here we choose tables URL
	spreadSheet, err := service.FetchSpreadsheet("1F2jHgop053dU8ITHOBLEAtMmgVgZYYYXOt2nWWIuR6M")
	CheckError(err)

	//here we choose table num
	sheet, err := spreadSheet.SheetByIndex(0)
	CheckError(err)

	sheet.Update(row, column, value)

	// Make sure call Synchronize to reflect the changes
	err = sheet.Synchronize()
	CheckError(err)
}

func GoogleSheetsCheckColumn() {
	data, err := os.ReadFile("local-abbey-420416-ecb7d8e898c0.json")
	CheckError(err)
	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	CheckError(err)
	client := conf.Client(context.TODO())

	service := spreadsheet.NewServiceWithClient(client)

	//here we choose tables URL
	spreadSheet, err := service.FetchSpreadsheet("1rk_i8y9_QNLLCLK2qYvUBEAMKcWkbM1WP2YNngeecMg")
	CheckError(err)

	//here we choose table num
	sheet, err := spreadSheet.SheetByIndex(0)
	CheckError(err)

	fmt.Println("sheet.Rows:", sheet.Rows)
	fmt.Println("sheet.Columns:", sheet.Columns)
	fmt.Println("sheet.Data:", sheet.Data)
	fmt.Println("sheet.Columns", sheet.Columns)

	// Make sure call Synchronize to reflect the changes
	err = sheet.Synchronize()
	CheckError(err)
}

func CheckError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
