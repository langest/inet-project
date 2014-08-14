package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var (
	response string
	filePath = "test/"
)

func main() {
	response = readFile("testPage.html")
	http.HandleFunc("/", handler) //Redirect all urls to handler function
	http.ListenAndServe("localhost:8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, response)
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
	}
	return strings.Join(lines, "")
}
