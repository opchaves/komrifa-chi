package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", helloWorld)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("could not listen on port 8080 %v", err)
	}
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name == "" {
		name = "World"
	}

	fmt.Fprintf(w, "Hello, %s!", name)
}
