package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"helloWorld/config"
	"helloWorld/database"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func fileUploadGoogleDrive(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	//file := &drive.File{Name: "Makefile", MimeType: "text/plain"}
	res := srv.Files.Export("Make", "text/plain")
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	fmt.Println(res)

	//file := &drive.File{Name: "Makefile", MimeType: "text/plain"}
	//res, err := srv.Files.Create(file).Do()
	//if err != nil {
	//	log.Fatalf("Unable to retrieve files: %v", err)
	//}
	//fmt.Println(res)

	//res, err := srv.Files.List().PageSize(10).
	//	Fields("nextPageToken, files(id, name)").Do()
	//fmt.Println("Files:", res)
	//if len(res.Files) == 0 {
	//	fmt.Println("No files found.")
	//} else {
	//	for _, i := range res.Files {
	//		fmt.Printf("%s (%s)\n", i.Name, i.Id)
	//	}
	//}
}

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "This is the root command",
	Long:  "It is the root command of cobra",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Root command is here.")
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

func fileUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "File Upload")
	file, handler, err := r.FormFile("myfile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %v\n", handler.Filename)
	fmt.Printf("File Size: %v\n", handler.Size)

	res := strings.Split(handler.Filename, ".")
	typeOfFile := res[len(res)-1]
	fileName := fmt.Sprintf("upload-*.%s", typeOfFile)
	tempFile, err := ioutil.TempFile("images", fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func httpServer() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This Is Our Golang Server")
	})
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.HandleFunc("/upload", fileUpload)
	r.HandleFunc("/upload-google", fileUploadGoogleDrive)
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
}
