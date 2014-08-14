package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "lol")
}

func main() {
	http.HandleFunc("/", handler) //Redirect all urls to handler function
	http.ListenAndServe("localhost:8080", nil)
}
