package main

import (
	"archiver/cmd/api/store"
	"fmt"
	"log"
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
		port:  portInt,
		store: store.InitStore(store.Azure | store.FileSystem | store.TempFileSystem), // could use some env variable config
	}

	app.start()
}

func (a *application) start() {
	mux := a.getRoutes()
	addr := fmt.Sprintf(":%d", a.port)

	log.Printf("Starting server at port %d\n", a.port)
	http.ListenAndServe(addr, mux)
}
