package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"helloWorld/database"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
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

var (
	oauthConfGl = &oauth2.Config{
		ClientID:     "1074893190578-h817e7n053f03iipha8ojosjl2kc5t2a.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8080/hello",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = "sanjay"
)

func googleOAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello")
	url := oauthConfGl.AuthCodeURL("sanjay")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) error {
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	fmt.Fprintf(w, "Content: %s\n", content)

	return json.NewEncoder(w).Encode(content)
}
func getUserInfo(state string, code string) ([]byte, error) {
	log.Println(state, code)
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := oauthConfGl.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	log.Println(token)
	token.Expiry = time.Unix(20000, 10)
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	fmt.Println("Body: ", response.Body)
	contents, err := ioutil.ReadAll(response.Body)
	fmt.Println("Content: ", string(contents))
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}

type Employee struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var user User
	user.Name = r.FormValue("username")
	insertUser(database.Get(), user.Name)
}
func GenerateJWT(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Minute * 2).Unix(),
		})
	var secret = []byte("sanjay")
	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Println("Error in generating token string", err.Error())
	}
	return tokenString
}
func validateJWT(dashboard http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		receivedToken, err := r.Cookie("token")
		if receivedToken == nil {
			fmt.Println("Received Token is nil")
			w.Write([]byte("Invalid JWT Token"))
			return
		}
		var key = []byte("sanjay")
		tokenString := receivedToken.Value
		fmt.Println("Received Token", tokenString)
		token, err := jwt.Parse(receivedToken.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("signing Method Unexpected")
			}
			return key, nil
		})
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims["username"], claims["exp"])
		} else {
			fmt.Println(err)
		}
		dashboard.ServeHTTP(w, r)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	var user User
	user.Name = r.FormValue("username")
	token := GenerateJWT(user.Name)
	fmt.Println("Token", token)

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(5 * time.Minute),
	})
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Dashboard")
	w.Write([]byte("Welcome To Dashboard"))
}
func httpServer() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This Is Our Golang Server")
		//handleGoogleCallback(w, r)
	})
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	//r.HandleFunc("/upload", service.FileUpload)
	//r.HandleFunc("/google-login", googleOAuth)
	////r.HandleFunc("/google-callback", googleOAuthCallback)
	//r.HandleFunc("/upload-google", service.FileUploadGoogleDrive)
	//r.HandleFunc("/users/{var}", handleUser)
	//r.HandleFunc("/register", register)
	//r.HandleFunc("/login", login)
	//r.HandleFunc("/dashboard", validateJWT(dashboard))
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
	//config.InitConfig()
	//config.InitConfiguration()
	//database.ConnectToDatabase()
	//Execute()
	//getAllUsers(database.Get())
	//config.InitializeLogger()
	httpServer()
	//database.CloseDatabase()
	//fmt.Println("Testing...")
}
