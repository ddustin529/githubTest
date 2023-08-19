package context

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func GetMySQLDB() *sql.DB {
	db, err := sql.Open("mysql", "dustin:dad1022@tcp(127.0.0.1:3306)/student")

	if err != nil {
		panic(err)
	}
	return db
}
