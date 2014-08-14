package main

import (
	"fmt"
	"net/http"
)

func handler() {

}

func main() {
	http.HandeFunc("/", handler) //Redirect all urls to handler function
	http.ListenAndServe("localhost:8080", nil)
}
