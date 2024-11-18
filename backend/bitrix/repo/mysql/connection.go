package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB //global variable

func InitDB(url string) (err error) {
	Db, err = sqlx.Connect("mysql", url)
	if err != nil {
		fmt.Println("not connected")
		return
	}
	err = Db.Ping()
	return
}
