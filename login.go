package main

import (
	"crypto/rand"
	"errors"
	"io"
	"log"

	"code.google.com/p/go.crypto/scrypt"
)

var (
	currentUser     = ""
	currentPassword = hashPassword("")
)

func NewUser(newUsername, newPassword string) {
	currentUser = newUsername
	currentPassword = newPassword
}

func Login(username, password string) error {
	if username == currentUser && password == currentPassword {
		return nil
	}
	return errors.New("login")
}
