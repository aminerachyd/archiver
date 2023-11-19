package store

import "log"

type Store interface {
	GetArchive(archiveName string) (archive, error)
	GetArchivesInfo() map[string]archiveMetadata
	PutArchive(archiveName string, payload []byte) error
	DeleteArchive(archiveName string) error
}

func InitStore() Store {
	store, err := InitMultiStore()

	if err != nil {
		log.Fatalf("Couldn't init multi store, error was [%s]\n", err)
	}
	return store
}
