package db

import (
	"database/sql"
	"fmt"
	"log"
)

func Conn(sgbd string) *sql.DB {
	var connStr string
	if sgbd == "mysql" {
		connStr = "urladmin:admin123@/urlshortnr"
	} else if sgbd == "postgres" {
		connStr = fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			"127.0.0.1", 5432, "arichard", "admin123", "mydb")
	} else {
		log.Fatal(connStr)
	}
	db, err := sql.Open(sgbd, connStr)
	if err != nil {
		panic(err)
	}
	return db
}

