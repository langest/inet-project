package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"errors"
	"log"
	"regexp"
	"time"
)

const (
	DATABASE_USER = "inetproj"
	DATABASE_PASS = "inetpass"
	DATABASE_NAME = "inetdb"
)

var (
	database *sql.DB
)

type noteInfo struct {
	title string
	note  string
}

func connectToDB() (*sql.DB, error) {
	database, err := sql.Open("mysql", DATABASE_USER+":"+DATABASE_PASS+"@/"+DATABASE_NAME)
	if err != nil {
		log.Println("Failed to connect to db")
		return nil, err
	}
	return database, nil
}

func closeDB() {
	database.Close()
}

func addUser(db *sql.DB, username, password string) (err error) {
	r := regexp.MustCompile(`[^a-zA-Z]+`)
	notOk := r.MatchString(username)
	if notOk {
		err = errors.New("username contains illegal characters")
		return
	}
	notOk = len(username) < 3
	if notOk {
		err = errors.New("username too short")
		return
	}
	notOk = len(password) < 10
	if notOk {
		err = errors.New("password too short")
		return
	}

	prepStmt, err := db.Prepare("SELECT * FROM users WHERE username = ?")
	if err != nil {
		return
	}
	rows, err := prepStmt.Query(username)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = errors.New("username already exists")
		return
	}

	prepStmt, err = db.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
	if err != nil {
		return
	}
	_, err = prepStmt.Exec(username, password)

	return
}

func checkPassword(db *sql.DB, username, password string) (ok bool, err error) {
	prepStmt, err := db.Prepare("SELECT * FROM users WHERE username = ? AND password = ?")
	if err != nil {
		return
	}
	rows, err := prepStmt.Query(username, password)
	if err != nil {
		return
	}
	//Check that we find exactly 1 user
	ok = rows.Next()
	ok = ok && !rows.Next()
	return
}

func getNotes(db *sql.DB, username string) (notes []noteInfo, err error) {
	prepStmt, err := db.Prepare("SELECT title, note FROM notes WHERE username = ? ORDER BY timestamp")
	if err != nil {
		return
	}

	rows, err := prepStmt.Query(username)
	if err != nil {
		return
	}

	notes = make([]noteInfo, 0)
	for rows.Next() {
		var ni noteInfo
		rows.Scan(&ni.title, &ni.note)
		notes = append(notes, ni)
	}
	return
}

func addNote(db *sql.DB, username, title, note string) (err error) {
	prepStmt, err := db.Prepare("INSERT INTO notes (username, title, note, timestamp) VALUES (?, ?, ?, ?)")
	if err != nil {
		return
	}
	_, err = prepStmt.Exec(username, title, note, time.Now())
	return
}

func removeNote(db *sql.DB, username, title string) (err error) {
	prepStmt, err := db.Prepare("DELETE FROM notes WHERE username = ? AND title = ?")
	if err != nil {
		return
	}
	_, err = prepStmt.Exec(username, title)
	return
}
