package main

import (
	"archiver/cmd/api/store"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

var DEFAULT_DESTINATIONS = []string{"azure, fs, tmpfs"}

func (a *application) healthHandler(w http.ResponseWriter, r *http.Request) {
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
	archiveName := r.URL.Path[len(V1_SINGLE_ARCHIVE):]
	if archiveName == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getSingleArchiveHandler(w, r, a.store, archiveName)
	case http.MethodPut:
		putSingleArchiveHandler(w, r, a.store, archiveName)
	case http.MethodDelete:
		deleteSingleArchiveHandler(w, r, a.store, archiveName)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func getSingleArchiveHandler(w http.ResponseWriter, r *http.Request, s store.Store, archiveName string) {
	archive, err := s.GetArchive(archiveName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(archive.Payload)
}

func putSingleArchiveHandler(w http.ResponseWriter, r *http.Request, s store.Store, archiveName string) {
	payload, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("error while PUT of [%s]. Error was [%s]\n", archiveName, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	destinations, present := r.Header["X-Storage-Destinations"]
	if !present {
		// If no destination specified, upload to all stores
		destinations = DEFAULT_DESTINATIONS
	}
	// Element at index 0 is a single string with header values separated by ","
	destinations = strings.Split(destinations[0], ",")
	log.Printf("Destinations to upload to [%v]", destinations)

	var wg sync.WaitGroup
	wg.Add(len(destinations))
	errCh := make(chan error, len(destinations))
	for _, dest := range destinations {
		dest = strings.TrimSpace(dest)
		storeType, err := store.Parse(dest)

		if err != nil {
			log.Printf("error while PUT of [%s]. Error was [%s]\n", archiveName, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		go func() {
			log.Printf("Uploading archive to store [%v]", storeType)
			errCh <- s.PutArchive(archiveName, payload, &storeType)
			wg.Done()
		}()
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			log.Printf("error while PUT of [%s]. Error was [%s]\n", archiveName, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func deleteSingleArchiveHandler(w http.ResponseWriter, r *http.Request, s store.Store, archiveName string) {
	err := s.DeleteArchive(archiveName)
	if err != nil {
		log.Printf("error while DELETE of [%s]. Error was [%s]\n", archiveName, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getArchivesInfoHandler(w http.ResponseWriter, r *http.Request, s store.Store) {
	payload := s.GetArchivesInfo()

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error while GET of archives infos. Error was [%s]\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	jsonPayload = append(jsonPayload, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonPayload)
}
