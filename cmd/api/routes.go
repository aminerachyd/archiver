package main

import (
	"net/http"
)

func (a *application) getRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/health", a.heatlhHandler)
	mux.HandleFunc("/v1/archives", a.archivesInfoHandler)
	mux.HandleFunc("/v1/archives/", a.singleArchiveHandler)

	return mux
}
