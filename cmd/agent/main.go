package main

import (
	"log"
	"net/http"

	"ospnet/internal/agent"
)

func main() {
	server := agent.NewServer()

	log.Println("agent starting on :9000")
	if err := http.ListenAndServe(":9000", server.Router()); err != nil {
		log.Fatal(err)
	}
}
