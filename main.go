package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"helloWorld/config"
	"helloWorld/database"
	"log"
	"net/http"
	"strconv"
)

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "This is the root command",
	Long:  "It is the root command of cobra",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Root command is here bro.")
	},
}

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "Add Two Numbers",
	Long:  "It will add two numbers.",
	Run: func(cmd *cobra.Command, args []string) {
		sum := 0
		for _, value := range args {
			num, err := strconv.Atoi(value)

			if err != nil {
				fmt.Println(err)
			}
			sum += num
		}
		fmt.Println(sum)
	},
}

func runMigrations() {
	driver, err := postgres.WithInstance(database.Get(), &postgres.Config{})
	if err != nil {
		fmt.Println(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://database/migration", "postgres", driver)
	if err != nil {
		fmt.Println("Migration Error")
	}

	result := m.Up()
	if result == nil {
		fmt.Println("After Migration Error", result)
	}
}

var migrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "It will run all the migrations",
	Long:  "It will run all the migrations from the CLI",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations()
	},
}

func init() {
	rootCmd.AddCommand(addCommand)
	rootCmd.AddCommand(migrateCommand)
}
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		return
	}
}

type User struct {
	name string
}

func getAllUsers(db *sql.DB) []string {
	var users []string
	sqlStatement := `SELECT * FROM users;`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Println("Error In Database Insert", err)
	}
	user := new(User)
	for rows.Next() {
		rows.Scan(&user.name)
		users = append(users, user.name)
	}
	return users
}

func httpServer() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This Is Our Golang Server")
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		users := getAllUsers(database.Get())
		json.NewEncoder(w).Encode(users)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
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
	Execute()
	getAllUsers(database.Get())
	config.InitializeLogger()
	insertIntoDatabase(database.Get(), "Sanjay Singh")
	httpServer()
	database.CloseDatabase()
}
