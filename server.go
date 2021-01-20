package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

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
	check(db.Ping())
	defer tx.Rollback()
	for data := range ["https://google.fr", "https://youtube.com"] {
		insert(tx, data)
	}

	// var hashedURL string

	data, err := tx.Prepare("SELECT url FROM storage WHERE url = ?")
	check(err)

	defer data.Close()
	results, _ := data.Exec("https://github.com")
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(results)
	// db.SetConnMaxLifetime(time.Minute * 3)
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(10)

	// r := mux.NewRouter()
	// r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)
	// srv := &http.Server{
	// 	Handler:      r,
	// 	Addr:         "127.0.0.1:8000",
	// 	WriteTimeout: 15 * time.Second,
	// 	ReadTimeout:  15 * time.Second,
	// }
	// log.Fatal(srv.ListenAndServe())
}

func insert(db *sql.Tx, data string) {
	insertCursor, err := db.Prepare("INSERT INTO storage VALUES( url, ? )")
	check(err)
	insertCursor.Exec(data)
	if err := db.Commit(); err != nil {
		log.Fatal(err)
	} else {
		log.Println(fmt.Sprintf("%s inserted in db", data))
	}
}
