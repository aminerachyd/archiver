package store

import (
	"log"
	"sync"
)

type MultiStore struct {
	stores []Store
}

func InitMultiStore(flags storageType) (Store, error) {
	stores := []Store{}

	if flags&Azure != 0 {
		log.Println("Enabling Azure store")
		azureStore, err := InitAzureStore()
		if err != nil {
			return nil, err
		}
		stores = append(stores, azureStore)
	}

	if flags&FileSystem != 0 {
		log.Println("Enabling file system store")
		fsStore, err := InitFileSystemStore()
		if err != nil {
			return nil, err
		}
		stores = append(stores, fsStore)
	}

	if flags&TempFileSystem != 0 {
		log.Println("Enabling temp file system store")
		tmpFsStore, err := InitTempFileSystemStore()
		if err != nil {
			return nil, err
		}
		stores = append(stores, tmpFsStore)
	}

	multiStore := MultiStore{
		stores: stores,
	}

	return multiStore, nil
}

func (s MultiStore) GetArchive(archiveName string) (*archive, error) {
	archiveCh := make(chan *archive)
	errCh := make(chan error, 1)

	for _, store := range s.stores {
		storeClone := store
		go func() {
			archive, err := storeClone.GetArchive(archiveName)
			if err != nil {
				errCh <- err
			} else {
				archiveCh <- archive
			}
		}()
	}

	var err error

	for i := len(s.stores); i > 0; {
		select {
		case <-archiveCh:
			return <-archiveCh, nil
		case <-errCh:
			err = <-errCh
			i--
		}
	}

	return nil, err
}

func (s MultiStore) GetArchivesInfo() map[string]archiveMetadata {
	log.Printf("Fetching archives infos from stores")
	archivesInfoCh := make(chan map[string]archiveMetadata, len(s.stores))
	var wg sync.WaitGroup

	wg.Add(len(s.stores))
	for _, store := range s.stores {
		storeCopy := store
		go func() {
			archivesInfoCh <- storeCopy.GetArchivesInfo()
			wg.Done()
		}()
	}
	wg.Wait()
	close(archivesInfoCh)
	resultArchivesInfo := map[string]archiveMetadata{}

	for m := range archivesInfoCh {
		resultArchivesInfo = merge(m, resultArchivesInfo)
	}

	return resultArchivesInfo
}

func (s MultiStore) PutArchive(archiveName string, payload []byte, dest *storageType) error {
	errCh := make(chan error, len(s.stores))
	var wg sync.WaitGroup

	wg.Add(len(s.stores))
	for _, store := range s.stores {
		storeCopy := store
		go func() {
			errCh <- storeCopy.PutArchive(archiveName, payload, dest)
			wg.Done()
		}()
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s MultiStore) DeleteArchive(archiveName string) error {
	errCh := make(chan error, len(s.stores))
	var wg sync.WaitGroup

	wg.Add(len(s.stores))
	for _, store := range s.stores {
		storeCopy := store
		go func() {
			errCh <- storeCopy.DeleteArchive(archiveName)
			wg.Done()
		}()
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}
