package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// ArticleHandler is an API endpoint
func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
	fmt.Fprintf(w, "ID: %v\n", vars["id"])
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	db, err := sql.Open("mysql", "urladmin:admin123@/urlshortnr")
	check(err)
	defer db.Close()

	tx, err := db.Begin()
	check(err)
	defer tx.Rollback()

	check(db.Ping())
	insert_url(tx, "https://google.fr")

	// var hashedURL string

	rows, err := db.Query("SELECT * FROM storage")
	check(err)
	columns, err := rows.Columns()
	check(err)
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------------")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	serve()
}

func insert_url(db *sql.Tx, data string) {
	insertCursor, err := db.Prepare("INSERT INTO storage VALUES( url, ? )")
	check(err)
	insertCursor.Exec(data)
	if err := db.Commit(); err != nil {
		log.Fatal(err)
	} else {
		log.Println(fmt.Sprintf("%s inserted in db", data))
	}
}

func serve() {
	r := mux.NewRouter()
	r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
