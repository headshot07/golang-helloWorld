package main

import (
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
)
import "helloWorld/logger"

type Database struct {
	User     string
	DbName   string
	Host     string
	Password string
	Port     int
}

var name = "Sanjay"

func httpServer() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This Is Our Golang Server")
	})
}
func pingRequest() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
func configReadDatabase() Database {
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Println("yamlFile Error", err)
	}
	var result Database
	data := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(yamlFile), &data)
	database := data["DATABASE"]
	for key, element := range database.(map[string]interface{}) {
		if key == "HOST" {
			result.Host = element.(string)
		} else if key == "PASSWORD" {
			result.Password = element.(string)
		} else if key == "DB_NAME" {
			result.DbName = element.(string)
		} else if key == "PORT" {
			result.Port = element.(int)
		} else if key == "USER" {
			result.User = element.(string)
		}
	}
	return result
}

func insertIntoDatabase(db *sql.DB, userName string) {
	sqlStatement := `INSERT INTO users (name) VALUES ($1)`
	result, err := db.Exec(sqlStatement, userName)
	fmt.Println(result)
	if err != nil {
		fmt.Println("Error In Database Insert")
	}
}

func main() {
	var database Database
	httpServer()
	pingRequest()
	database = configReadDatabase()
	fmt.Println("Database Credentials", database)
	postgresConnect := fmt.Sprintf("user=%s dbname=%s sslmode=disable", database.User, database.DbName)
	db, err := sql.Open("postgres", postgresConnect)

	if err != nil {
		fmt.Println("Error")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Error In Database Ping")
	}
	insertIntoDatabase(db, "Sanjay")
	logger.SetLogger()
	logger.InfoLogger.Println("Info Logger")
	logger.WarningLogger.Println("Warning Logger")
	logger.ErrorLogger.Println("Error Logger")
	log.Println("Log Testing Again")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
