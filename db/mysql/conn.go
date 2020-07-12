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

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range scanArgs {
		scanArgs[i] = &values[i]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			fmt.Printf("scan result failed: %v\n", err)
		}
		for j, val := range values {
			if val != nil {
				record[columns[j]] = val
			}
		}
		records = append(records, record)
	}
	return records
}
