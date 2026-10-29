package main

import (
	"log"
	"net/http"

	handler "github.com/emeinardus/daybrief/api"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", handler.Handler)

	log.Println("Daybrief runs on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
