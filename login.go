package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"

	"code.google.com/p/go.crypto/scrypt"
)

var (
	currentUser     = ""
	currentPassword = hashPassword("")
)

const (
	PW_SALT_BYTES = 32
	PW_HASH_BYTES = 64
)

func hashPassword(password string) string {
	salt := make([]byte, PW_SALT_BYTES)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		log.Fatal(err)
	}

	hash, err := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, PW_HASH_BYTES)
	if err != nil {
		log.Fatal(err)
	}

	return hash
}

func NewUser(newUsername, newPassword string) {
	currentUser = newUsername
	currentpassword = hashPassword(newPassword)
}

func Login(username, password string) (string, error) {
	if username == currentUser && hashPassword(password) == currentPassword {
		return username, nil
	}
	return nil, errors.New("login")
}
