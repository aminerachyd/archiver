package main

import (
	"archiver/cmd/api/store"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type application struct {
	port  int
	store store.Store
}

func main() {
	port := os.Getenv("PORT")

	portInt, err := strconv.Atoi(port)
	if err != nil {
		portInt = 8080
	}

	app := application{
		port: portInt,
		// TODO Make this configurable, different types of stores ?
		store: store.InitAzureStore(),
	}

	app.start()
}

func (a *application) start() {
	mux := a.getRoutes()
	addr := fmt.Sprintf(":%d", a.port)

	fmt.Printf("Starting server at port %d\n", a.port)
	http.ListenAndServe(addr, mux)
}
