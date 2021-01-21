package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"./db"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/kare/base62"
	_ "github.com/lib/pq"
)

func encodeRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	rawURL := vars["rawURL"]
	encodedURL, _ := base62.Decode(rawURL)
	cursor := db.Conn("postgres")
	check(cursor.Ping())
	_, err := cursor.Query(fmt.Sprintf("INSERT INTO web_url (URL) VALUES ( '%s' );", encodedURL))
	check(err)
	cursor.Close()
}

func decodeRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	encodedURL, _ := strconv.ParseInt(vars["encodedURL"], 62, 64)
	rawURL := base62.Encode(encodedURL)
	http.Redirect(w, r, rawURL, http.StatusFound)
	log.Println(encodedURL)
	cursor := db.Conn("postgres")
	check(cursor.Ping())
	rows, err := cursor.Query("SELECT * FROM web_url;")
	check(err)
	results := fetch(rows)
	for _, result := range results {
		log.Println(fmt.Sprintf("%v : %s", result.id, result.url))
	}
	_ = cursor.Close()
}

func main() {
	serve()
}

func insertURL(tx *sql.Tx, data string) {
	stmt, err := tx.Prepare("INSERT INTO web_url VALUES (?);")
	if err != nil {
		log.Fatal(err)
	} else {
		_, err := stmt.Exec(data)
		check(err)
		log.Println(fmt.Sprintf("%s inserted in tx", data))
	}
}

func serve() {
	r := mux.NewRouter()
	r.HandleFunc("/{encodedURL}", decodeRoute).Methods("GET")
	r.HandleFunc("encode/{rawURL}", encodeRoute)
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}



type webURL struct {
	id  int
	url string
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
