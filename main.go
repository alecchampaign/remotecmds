package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received.")
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
