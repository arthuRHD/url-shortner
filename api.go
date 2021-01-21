package main

import (
	"database/sql"
	"fmt"
	"github.com/kare/base62"
	"log"
	"net/http"
	"strconv"
	"time"
	 "./db"

	"github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

// GetOriginalURL fetches the original URL for the given encoded(short) string
func (driver *DBClient) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	encoded := vars["encoded_string"]
	res := base62.Encode(encoded)
	log.Println(res)
	if encodedURL, err := strconv.ParseInt(encoded, 10, 64); err != nil {
		panic(err)
	} else {
		base62.Encode()
		rawURL, _ := base62.Decode(encodedURL)
		http.Redirect(w, r, rawURL, http.StatusFound)
	}
}

// GenerateShortURL adds URL to DB and gives back shortened string
func (driver *DBClient) GenerateShortURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	rawURL := vars[""]
	encodedInt, _ := base62.Decode(rawURL)
	rows, err := driver.db.Query(fmt.Sprintf("INSERT INTO web_url VALUES (%v, '%s')", encodedInt, rawURL))
	checkErr(err)
	fetch(rows)
	http.Redirect(w, r, rawURL, http.StatusFound)
}
		
func main() {
	cursor := db.Conn("postgres")
	dbClient := &DBClient{db: cursor}
	defer cursor.Close()
	r := mux.NewRouter()
	r.HandleFunc("/v1/{encoded_string:[a-zA-Z0-9]*}", dbClient.GetOriginalURL).Methods("GET")
	r.HandleFunc("/v1/new", dbClient.GenerateShortURL).Methods("POST")
	srv := &http.Server{
		Handler: r,
		Addr:"127.0.0.1:8000",
		WriteTimeout:15* time.Second,
		ReadTimeout:15* time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
		
func checkErr(e error) {
	if e !=nil {
		panic(e)
	}
}
// DB stores the database session information. Needs to be initialized once
type DBClient struct {
	db *sql.DB
}

// Model the record struct
type Record struct {
	ID  int `json:"id"`
	URL string `json:"url"`
}

func fetch(rows *sql.Rows) []Record {
	var results []Record
	if rows != nil {
		for rows.Next() {
			var (
				id  int
				url string
			)
			if err := rows.Scan(&id, &url); err != nil {
				panic(err)
			}
			results = append(results, Record{id, url})
		}
	} else {
		log.Fatal(rows)
	}
	return results
}