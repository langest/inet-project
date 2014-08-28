package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	filePath = "test/"
	store    = sessions.NewCookieStore([]byte("something-very-secret"))
)

func main() {
	http.HandleFunc("/", handler) //Redirect all urls to handler function
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/register", handleRegister)
	err := http.ListenAndServeTLS("localhost:8080", filePath+"cert.pem", filePath+"key.pem", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, readFile("testPage.html"))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, readFile("login.html"))
		fmt.Println("No post request")
	} else {
		username := r.FormValue("username")
		password := r.FormValue("password")
		fmt.Printf("username: %s, password: %s\n", username, password)
		fmt.Fprintf(w, readFile("login.html"))
		fmt.Println("Post request")
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, readFile("register.html"))
}

func readFile(fileName string) string {
	file, err := os.Open(filePath + fileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	scanner := bufio.NewScanner(file)

	lines := make([]string, 1024)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		lines = append(lines, "\n")
	}
	return strings.Join(lines, "")
}
