package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"errors"
	"regexp"
)

const (
	DATABASE_USER = "root"
	DATABASE_PASS = "pass"
	DATABASE_NAME = "intnet"
)

var (
	database *sql.DB
)

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
	ok, err := regexp.MatchString("[a-zA-Z0-9_]", username)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New("username contains illegal characters")
	}
	//TODO chack that username isnt used already
	prepStmt, err := db.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
	if err != nil {
		return
	}
	prepStmt.Exec(username, password)
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
	rows.Next()
	ok := rows.Next()
	return
}

func addNote(db *sql.DB, username, note string) (err error) {
	prepStmt, err := db.Prepare("INSERT INTO notes (username, note) VALUES (?, ?) ON DUPLICATE KEY UPDATE username = ?")
	if err != nil {
		return
	}
	_, err := prepStmt.Exec(username, note, username)
	return
}

func removeNote(db *sql.DB, username, note string) (err error) {
	prepStmt, err := db.Prepare("DELETE notes WHERE username = ? AND note = ?")
	if err != nil {
		return
	}
	_, err := prepStmt.Exec(username, note)
	return
}
