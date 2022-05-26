package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"helloWorld/config"
	"helloWorld/database"
	"helloWorld/service"
	"log"
	"net/http"
	"strconv"
)

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "This is the root command",
	Long:  "It is the root command of cobra",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Root command is here.")
	},
}

var migrateCommand = &cobra.Command{
	Use:   "migrate-up",
	Short: "It will run all the migrations",
	Long:  "It will run all the migrations from the CLI",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations()
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
	Name string `json:"name"`
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
		rows.Scan(&user.Name)
		users = append(users, user.Name)
	}
	return users
}

func insertUser(db *sql.DB, name string) {
	sqlStatement := `INSERT INTO users (name) VALUES ($1);`
	_, err := db.Exec(sqlStatement, name)
	if err != nil {
		fmt.Println("Error In Database Insert", err)
	}
}

func deleteUser(db *sql.DB, username string) {
	fmt.Println("User Name", username)
	sqlStatement := `DELETE FROM users WHERE name=($1);`
	_, err := db.Exec(sqlStatement, username)
	if err != nil {
		fmt.Println("Error In Database Delete", err)
	}
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	var param = mux.Vars(r)
	if r.Method == "GET" {
		users := getAllUsers(database.Get())
		json.NewEncoder(w).Encode(users)
		fmt.Println("Get Users Method")
	} else if r.Method == "POST" {
		user := new(User)
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			fmt.Println("Decode Error")
		}
		fmt.Println("Request Body", user.Name)
		insertUser(database.Get(), user.Name)
		w.Write([]byte("User Added Successfully"))
	} else if r.Method == "DELETE" {
		fmt.Println("Delete", param["var"])
		deleteUser(database.Get(), param["var"])
		w.Write([]byte("User Deleted Successfully"))
	}
}

func httpServer() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This Is Our Golang Server")
	})
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.HandleFunc("/upload", service.FileUpload)
	r.HandleFunc("/upload-google", service.FileUploadGoogleDrive)
	r.HandleFunc("/users/{var}", handleUser)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
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
	httpServer()
	database.CloseDatabase()
	fmt.Println("Testing...")
}
