package store

type Store interface {
	GetArchive(archiveName string) (archive, error)
	PutArchive(archiveName string, payload []byte) error
	GetArchivesInfo() []archiveMetadata
}
