package store

import (
	"fmt"
	"log"
	"os"
)

type FileSystemStore struct {
	archivesPath string
	storageType  storageType
}

func InitFileSystemStore() (Store, error) {
	archivesPath := os.Getenv("ARCHIVER_FILESYSTEM_PATH")

	return initWithPath(archivesPath, false)
}

func InitTempFileSystemStore() (Store, error) {
	archivesPath := "/tmp/archives"

	return initWithPath(archivesPath, true)
}

func initWithPath(path string, isTemp bool) (Store, error) {
	storageType := FileSystem
	if isTemp {
		storageType = TempFileSystem
	}

	if _, err := os.ReadDir(path); err != nil {
		err = os.Mkdir(path, 0777)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialize file system store of type [%s]. error was [%s]", storageType.toString(), err)
		}
	}

	store := FileSystemStore{
		archivesPath: path,
		storageType:  storageType,
	}

	return store, nil
}

func (s FileSystemStore) GetArchive(archiveName string) (*archive, error) {
	filePath := fmt.Sprintf("%s/%s", s.archivesPath, archiveName)
	payload, err := os.ReadFile(filePath)

	archive := archive{
		Payload: payload,
		metadata: archiveMetadata{
			SizeInBytes: int64(len(payload)),
			StoredIn:    []string{s.storageType.toString()},
		},
	}

	return &archive, err
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
					SizeInBytes: fileSize,
					StoredIn:    []string{s.storageType.toString()},
				}
				archivesMetadata[fileName] = metadata
			}
		}
	}
	return archivesMetadata
}

func (s FileSystemStore) PutArchive(archiveName string, payload []byte, dest *storageType) error {
	if dest == nil {
		return fmt.Errorf("no destination specified for Filesystem [%s] store", s.storageType.toString())
	}

	if *dest == FileSystem || *dest == TempFileSystem {
		filePath := fmt.Sprintf("%s/%s", s.archivesPath, archiveName)
		err := os.WriteFile(filePath, payload, 0666)
		return err
	}

	log.Printf("wrong destination specified for Filesystem [%s] store [%v], skipping upload", s.storageType.toString(), dest.toString())
	return nil
}

func (s FileSystemStore) DeleteArchive(archiveName string) error {
	filePath := fmt.Sprintf("%s/%s", s.archivesPath, archiveName)
	err := os.Remove(filePath)
	return err
}
