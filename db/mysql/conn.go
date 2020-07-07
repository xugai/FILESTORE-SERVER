package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:xuqiubing@tcp(127.0.0.1)/fileserver?charset=utf8")
	err := db.Ping()
	if err != nil {
		fmt.Printf("Failed to ping succeed with mysql: %v\n", err)
		os.Exit(1)
	}
}

func GetDBConnection() *sql.DB {
	return db
}
