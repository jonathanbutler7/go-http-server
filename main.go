package main

import (
	"net/http"

	"example.com/m/api"
)

func main() {
	srv := api.NewServer()
	http.ListenAndServe(":8080", srv)
}
