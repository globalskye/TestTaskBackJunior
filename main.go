package main

import (
	"jwttask/config"
	"jwttask/handlers"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/register", handlers.Register) //method get

	http.HandleFunc("/refresh", handlers.Refresh) //method post

	log.Fatal(http.ListenAndServe(config.Conf.ListenAddress, nil))

}
