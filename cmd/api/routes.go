package main

import (
	"net/http"
)

const V1_HEALTH = "/v1/health"
const V1_ARCHIVES = "/v1/archives"
const V1_SINGLE_ARCHIVE = "/v1/archives/"

func (a *application) getRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc(V1_HEALTH, a.healthHandler)
	mux.HandleFunc(V1_ARCHIVES, a.archivesInfoHandler)
	mux.HandleFunc(V1_SINGLE_ARCHIVE, a.singleArchiveHandler)

	return mux
}
