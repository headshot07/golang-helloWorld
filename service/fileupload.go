package service

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	fmt.Printf("Go to the following url %v", authURL)

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
	fmt.Printf("Saving credentials file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func returnContentType(typeOfFile string) string {
	if typeOfFile == "pdf" {
		return "application/pdf"
	} else if typeOfFile == "txt" {
		return "text/plain"
	} else if typeOfFile == "jpg" || typeOfFile == "jpeg" {
		return "image/jpeg"
	} else if typeOfFile == "zip" {
		return "application/zip"
	} else {
		return "text/plain"
	}
}
func FileUploadGoogleDrive(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	file, handler, err := r.FormFile("myfile")

	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	result := strings.Split(handler.Filename, ".")
	typeOfFile := result[len(result)-1]

	contentType := returnContentType(typeOfFile)
	uploadFile := &drive.File{Name: handler.Filename}

	if err != nil {
		log.Fatalln(err)
	}
	_, _ = srv.Files.Create(uploadFile).Media(file, googleapi.ContentType(contentType)).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	//file := &drive.File{Name: "Makefile", MimeType: "text/plain" }
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

func FileUpload(w http.ResponseWriter, r *http.Request) {
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
