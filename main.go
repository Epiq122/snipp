package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from the snipp"))
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	log.Print("Server started on port 8080")

	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
