package store

import "log"

type Store interface {
	GetArchive(archiveName string) (archive, error)
	GetArchivesInfo() []archiveMetadata
	PutArchive(archiveName string, payload []byte) error
	DeleteArchive(archiveName string) error
}

func InitStore() Store {
	// TODO Make this configurable, different types of stores ?
	store, err := InitAzureStore()

	if err != nil {
		log.Fatalf("Couldn't init Azure store, error was [%s]\n", err)
	}
	return store
}
