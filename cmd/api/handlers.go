package main

import (
	"archiver/cmd/api/store"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (a *application) heatlhHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "UP")
}

func (a *application) archivesInfoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getArchivesInfoHandler(w, r, a.store)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (a *application) singleArchiveHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getSingleArchiveHandler(w, r, a.store)
	case http.MethodPut:
		putSingleArchiveHandler(w, r, a.store)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func getSingleArchiveHandler(w http.ResponseWriter, r *http.Request, s store.Store) {
	archiveName := r.URL.Path[len("/v1/archives/"):]
	archive, err := s.GetArchive(archiveName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(archive.Payload)
}

func putSingleArchiveHandler(w http.ResponseWriter, r *http.Request, s store.Store) {
	archiveName := r.URL.Path[len("/v1/archives/"):]
	payload, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = s.PutArchive(archiveName, payload)
	if err != nil {
		log.Printf("Got error %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func getArchivesInfoHandler(w http.ResponseWriter, r *http.Request, s store.Store) {
	payload := s.GetArchivesInfo()

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Got error %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	jsonPayload = append(jsonPayload, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonPayload)
}
