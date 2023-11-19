package store

import (
	"fmt"
	"log"
	"os"
)

type FileSystemStore struct {
	archivesPath string
}

func InitFileSystemStore() (Store, error) {
	archivesPath := os.Getenv("ARCHIVER_FILESYSTEM_PATH")

	if _, err := os.ReadDir(archivesPath); err != nil {
		err = os.Mkdir(archivesPath, 0777)
		if err != nil {
			return nil, err
		}
	}

	store := FileSystemStore{
		archivesPath: archivesPath,
	}

	return store, nil
}

func (s FileSystemStore) GetArchive(archiveName string) (archive, error) {
	filePath := fmt.Sprintf("%s/%s", s.archivesPath, archiveName)
	payload, err := os.ReadFile(filePath)

	archive := archive{
		metadata: archiveMetadata{
			Name:        archiveName,
			SizeInBytes: int64(len(payload)),
		},
		Payload: payload,
	}

	return archive, err
}

func (s FileSystemStore) GetArchivesInfo() map[string]archiveMetadata {
	archivesMetadata := map[string]archiveMetadata{}
	dirEntries, err := os.ReadDir(s.archivesPath)
	if err != nil {
		log.Printf("Got error %s\n", err)
	} else {
		for _, dir := range dirEntries {
			if dir.Type().IsRegular() {
				fileName := dir.Name()
				info, err := dir.Info()
				if err != nil {
					log.Printf("Got error %s\n", err)
					continue
				}
				fileSize := info.Size()

				metadata := archiveMetadata{
					Name:        fileName,
					SizeInBytes: fileSize,
				}
				archivesMetadata[fileName] = metadata
			}
		}
	}
	return archivesMetadata
}

func (s FileSystemStore) PutArchive(archiveName string, payload []byte) error {
	filePath := fmt.Sprintf("%s/%s", s.archivesPath, archiveName)
	err := os.WriteFile(filePath, payload, 0666)
	return err
}

func (s FileSystemStore) DeleteArchive(archiveName string) error {
	filePath := fmt.Sprintf("%s/%s", s.archivesPath, archiveName)
	err := os.Remove(filePath)
	return err
}
