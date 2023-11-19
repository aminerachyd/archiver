package store

import "maps"

// TODO Manage
// - Concurrent case, goroutines go brr
// - Single store used in multistore case
type MultiStore struct {
	stores []Store
}

func InitMultiStore() (Store, error) {
	stores := []Store{}

	azureStore, err := InitAzureStore()
	if err != nil {
		return nil, err
	}
	stores = append(stores, azureStore)

	fsStore, err := InitFileSystemStore()
	if err != nil {
		return nil, err
	}
	stores = append(stores, fsStore)

	multiStore := MultiStore{
		stores: stores,
	}

	return multiStore, nil
}

func (s MultiStore) GetArchive(archiveName string) (archive, error) {
	var resultArchive archive
	var err error
	found := false
	for _, store := range s.stores {
		resultArchive, err = store.GetArchive(archiveName)
		if err == nil {
			found = true
			break
		}
	}

	if found {
		return resultArchive, nil
	} else {
		return resultArchive, err
	}
}

func (s MultiStore) GetArchivesInfo() map[string]archiveMetadata {
	resultArchivesInfo := map[string]archiveMetadata{}
	for _, store := range s.stores {
		maps.Copy(resultArchivesInfo, store.GetArchivesInfo())
	}

	return resultArchivesInfo
}

func (s MultiStore) PutArchive(archiveName string, payload []byte) error {
	for _, store := range s.stores {
		err := store.PutArchive(archiveName, payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s MultiStore) DeleteArchive(archiveName string) error {
	for _, store := range s.stores {
		err := store.DeleteArchive(archiveName)
		if err != nil {
			return err
		}
	}

	return nil
}
