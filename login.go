package main

import (
	"errors"
)

var (
	currentUser     = ""
	currentPassword = ""
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
