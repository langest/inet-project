package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
)

var (
	filePath     = "test/"
	sessionStore = sessions.NewCookieStore([]byte("something-very-secret"))
)

func main() {
	http.HandleFunc("/", handleIndex) //Redirect all urls to handler function
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/register", handleRegister)
	err := http.ListenAndServeTLS("localhost:8080", filePath+"cert.pem", filePath+"key.pem", context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, readFile("index.html"))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, readFile("login.html"))
	} else {
		username := r.FormValue("username")
		password := r.FormValue("password")
		err := Login(username, password)
		if err == nil {
			fmt.Println("logged in successfully")
			//TODO show successful login page and redirect to home or something
		} //TODO else show unsuccessful and show login again
		fmt.Fprintf(w, readFile("login.html"))
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, readFile("register.html"))
	} else {
		username := r.FormValue("username")
		password := r.FormValue("password")
		NewUser(username, password)
		handleIndex(w, r)
	}
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
