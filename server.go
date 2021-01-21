package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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
		log.Fatal(e)
	}
}

func main() {
	db := pgDB()
	defer db.Close()

	check(db.Ping())
	tx, err := db.Begin()
	check(err)
	defer tx.Rollback()

	fetchAll(db, "SELECT * FROM web_url;")
	log.Println("ça marche")
	insertURL(tx, "https://.fr")
	time.Sleep(time.Second * 5)
	fetchAll(db, "SELECT * FROM web_url;")
}

func insertURL(db *sql.Tx, data string) {
	log.Println("ça marche")
	query := "INSERT INTO web_url (URL) VALUES ( ? );"
	log.Println(query)
	stmt, err := db.Prepare(query)
	log.Println("ça marche")
	if err != nil {
		log.Fatal(err)
	} else {
		defer stmt.Close()
		log.Println("ça marche")
		_, err := stmt.Exec(data)
		check(err)
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

func fetchAll(db *sql.DB, query string) {
	var (
		url string
		id  int
	)
	rows, err := db.Query(query)
	check(err)
	defer rows.Close()
	for rows.Next() {
		check(rows.Scan(&id, &url))
		log.Println(id, url)
	}
}

func mysqlDB() *sql.DB {
	db, err := sql.Open("mysql", "urladmin:admin123@/urlshortnr")
	check(err)
	return db
}

func pgDB() *sql.DB {
	const (
		host     = "127.0.0.1"
		port     = 5432
		user     = "bberenger"
		password = "cesi2021"
		dbname   = "mydb"
	)
	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connectionString)
	check(err)
	return db
}
