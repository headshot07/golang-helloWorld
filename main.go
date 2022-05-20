package main

import (
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"helloWorld/config"
	"helloWorld/database"
	"log"
	"net/http"
)

func httpServer() {
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This Is Our Golang Server")
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func insertIntoDatabase(db *sql.DB, userName string) {
	sqlStatement := `INSERT INTO users (name) VALUES ($1)`
	_, err := db.Exec(sqlStatement, userName)
	if err != nil {
		fmt.Println("Error In Database Insert", err)
	}
}
func Add(x, y int) (res int) {
	return x + y
}

func compareString(str1 string, str2 string) bool {
	if str1 == str2 {
		return true
	}
	return false
}

func Multiply(a, b int) int {
	return a * b
}

func main() {
	config.InitConfig()
	config.InitConfiguration()
	database.ConnectToDatabase()
	config.InitializeLogger()
	insertIntoDatabase(database.Get(), "Sanjay Singh")
	database.CloseDatabase()
	httpServer()
}
