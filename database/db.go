package database

import (
	"database/sql"
	"fmt"
	"helloWorld/config"
)

var db *sql.DB

func ConnectToDatabase() {
	configuration := config.GetConfig()
	postgresConnect := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", configuration.Database.User, configuration.Database.Password, configuration.Database.DB_Name)
	database, err := sql.Open("postgres", postgresConnect)
	if err != nil {
		fmt.Println("Error")
	}
	db = database
	err = database.Ping()
	if err != nil {
		fmt.Println("Error In Database Ping")
	}
}
func CloseDatabase() {
	defer db.Close()
	db.Close()
}

func Get() *sql.DB {
	return db
}
