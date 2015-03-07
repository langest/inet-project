package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"os"
	"strings"
	
)

const (
	LISTENPORT = "8080"
)

var (
	filePath     = "test/"
	sessionStore = sessions.NewCookieStore([]byte("something-very-secret"))
	db           *sql.DB
)

func main() {
	var err error
	db, err = connectToDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer closeDB()

	log.Println("starting server...")
	http.HandleFunc("/", handleIndex) //Redirect all urls to handler function
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/loggedinpage", handleLoggedInPage)
	http.HandleFunc("/notes", handleNotes)
	err = http.ListenAndServeTLS("localhost:"+LISTENPORT, filePath+"cert.pem", filePath+"key.pem", context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, buildWebpage("", readFile("index.html")))
}

func handleLoggedInPage(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "login")
	if err != nil {
		//TODO
	}
	username, ok := session.Values["username"]
	if !ok {
		http.Redirect(w, r, "../login", http.StatusFound)
		return
	}
	fmt.Fprintf(w, buildWebpage("", readFile("loggedinpage.html")), username)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "login")
	ok := false
	if err == nil {
		_, ok = session.Values["username"]
	}
	if ok {
		log.Println("Already logged in as the user", session.Values["username"])
		http.Redirect(w, r, "../loggedinpage", http.StatusFound)
		return
	}

	//if it was not a post request
	//print the login page
	if r.Method != "POST" {
		fmt.Fprintf(w, buildWebpage(readFile("login_head.html"), readFile("login.html")))

	} else { //else try to login
		username := r.FormValue("username")
		password := r.FormValue("password")
		ok, err := checkPassword(db, username, password)
		if err != nil {
			log.Fatal(err)
		} else if ok {
			log.Println("logged in successfully")
			session.Values["username"] = username
			//TODO show successful login page and redirect to home or something

			session.Save(r, w)

		} //TODO else show unsuccessful and show login again
		fmt.Fprintf(w, buildWebpage(readFile("login_head.html"),readFile("login.html")))
	}
}

func handleNotes(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "login")
	ok := false
	if err == nil {
		_, ok = session.Values["username"]
	}
	if !ok {
		http.Redirect(w, r, "../login", http.StatusFound)
	}
	fmt.Fprintf(w, buildWebpage("", readFile("notes.html")))
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, buildWebpage(readFile("Crypto.html"), readFile("register.html")))
	} else {
		err := addUser(db, r.FormValue("username"), r.FormValue("password"))
		if err != nil {
			log.Fatal(err)
		}
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

func buildWebpage(header string, body string ) string {
    return readFile("start.html") + header + readFile("middle.html") + body + readFile("end.html");
}

