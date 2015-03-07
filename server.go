package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

const (
	LISTENPORT = "8080"
)

var (
	filePath     = "test/"
	sessionStore = sessions.NewCookieStore([]byte("something-very-secret"))
)

func main() {
	log.Println("starting server...")
	http.HandleFunc("/", handleIndex) //Redirect all urls to handler function
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/loggedinpage", handleLoggedInPage)
	err := http.ListenAndServeTLS("localhost:"+LISTENPORT, filePath+"cert.pem", filePath+"key.pem", context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, readFile("index.html"))
}

func handleLoggedInPage(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "login")
	if err != nil {
		//TODO
	}
	username, ok := session.Values["username"]
	if !ok {
		fmt.Fprintf(w, readFile("login.html"))
		return
	}
	fmt.Fprintf(w, readFile("loggedinpage.html"), username)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "login")
	if err != nil {
		//TODO
	}
	_, ok := session.Values["username"]
	if ok {
		log.Println("Already logged in as the user", session.Values["username"])
		fmt.Fprintf(w, readFile("login.html"))
		return
	}

	//if it was not a post request
	//print the login page
	if r.Method != "POST" {
		fmt.Fprintf(w, readFile("login.html"))

	} else { //else try to login
		username := r.FormValue("username")
		password := r.FormValue("password")
		err := Login(username, password)
		if err == nil {
			log.Println("logged in successfully")
			session.Values["username"] = username
			//TODO show successful login page and redirect to home or something
			session.Save(r, w)

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
