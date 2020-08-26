package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Person Struct
type Person struct {
	ID   string `json:"idEntity"`
	Name string `json:"Name"`
}

func main() {
	fmt.Println("Running server...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// set up database connection
			db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/ttscorporate_dev")
			defer db.Close()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = db.Ping()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			db.SetConnMaxLifetime(time.Minute * 3)
			db.SetMaxOpenConns(10)
			db.SetMaxIdleConns(10)

			// query the database
			results, err := db.Query("SELECT idEntity, CONCAT_WS(' ', FirstName, LastName) as Name FROM person LIMIT 10")
			defer results.Close()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// parse db rows to struct
			var persons []Person
			for results.Next() {
				var person Person
				err := results.Scan(&person.ID, &person.Name)
				if err != nil {
					if err != sql.ErrNoRows {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				} else {
					persons = append(persons, person)
				}
			}

			// convert to json
			json, err := json.Marshal(persons)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// send response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(json)
		}
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
