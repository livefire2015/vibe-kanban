package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	store := NewStore()
	mux := http.NewServeMux()
	RegisterHandlers(mux, store)

	mux.Handle("/", http.FileServer(http.Dir("frontend/dist")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
