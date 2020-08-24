package conn

import (
	"FILESTORE-SERVER/service/dbproxy/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", config.MySQLSource)
	db.SetMaxOpenConns(100)
	err := db.Ping()
	if err != nil {
		log.Printf("Failed to connect db: %v\n", err)
		os.Exit(1)
	}
}

func DBConn() *sql.DB {
	return db
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			panic(err)
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
