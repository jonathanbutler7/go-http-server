package main

import (
	"example.com/m/api"
	"log"
	"net/http"
)

func main() {
	srv := api.NewServer()
	err := http.ListenAndServe(":8080", srv)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
		// you should plan on always handling any errors in go unless you have a very good reason.
		// If your app can't start serving the http server on a port, that is a critical startup failure and you should
		// exit with a non-zero status code.
		// Alternatively, you can use `panic(err)` to crash the program. Fine, but go panics are fairly rare and instill panic in me haha
	}
}
