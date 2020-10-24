package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello!"))
	return
}

func main() {
	fmt.Println("Starting server on port 7777")

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRoute)

	log.Fatal(http.ListenAndServe(":7777", mux))
}
