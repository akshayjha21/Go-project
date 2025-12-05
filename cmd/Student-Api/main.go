package main

import (
	"fmt"
	"log"
	"net/http"

	Config "github.com/akshayjha21/Student-Api/internal/config"
)

func main() {
	//TODO - load config

	cfg := Config.MustLoad()

	//TODO - database setup

	//TODO - setup route
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to the student api"))
	})

	//TODO - setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	fmt.Println("Server has started ")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to connect to the server")
	}
}
