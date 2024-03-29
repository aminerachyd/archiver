package store

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type AzureStore struct {
	ctx        context.Context
	Client     *azblob.Client
	container  string
	serviceURL string
	SASToken   string
}

func InitAzureStore() (Store, error) {
	ctx := context.Background()

	// TODO Better credentials/envs management ?
	serviceURL := os.Getenv("AZURE_SERVICE_URL")
	SASToken := os.Getenv("AZURE_SAS_TOKEN")
	containerName := os.Getenv("AZURE_ARCHIVES_CONTAINER")
	// Envs needed for this:
	// - AZURE_TENANT_ID
	// - AZURE_CLIENT_ID
	// - AZURE_CLIENT_SECRET
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialize Azure store. error was [%s]", err)
	}

	client, err := azblob.NewClient(serviceURL, credential, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialize Azure system store. error was [%s]", err)
	}

	store := AzureStore{
		ctx:        ctx,
		container:  containerName,
		Client:     client,
		serviceURL: serviceURL,
		SASToken:   SASToken,
	}

	return store, nil
}

func (s AzureStore) GetArchive(archiveName string) (*archive, error) {
	get, err := s.Client.DownloadStream(s.ctx, s.container, archiveName, nil)
	if err != nil {
		log.Printf("Got error %s\n", err)
		return nil, err
	}
	log.Printf("Fetched archive: %s\n", archiveName)

	data := bytes.Buffer{}
	retryReader := get.NewRetryReader(s.ctx, nil)
	_, err = data.ReadFrom(retryReader)
	if err != nil {
		return nil, err
	}

	archive := archive{
		Payload: data.Bytes(),
		metadata: archiveMetadata{
			StoredIn:    []string{Azure.toString()},
			SizeInBytes: int64(len(data.Bytes())),
		},
	}

	return &archive, nil
}

func (s AzureStore) GetArchivesInfo() map[string]archiveMetadata {
	archivesInfo := map[string]archiveMetadata{}

	pager := s.Client.NewListBlobsFlatPager(s.container, &azblob.ListBlobsFlatOptions{
		Include: azblob.ListBlobsInclude{Snapshots: true, Versions: true, Metadata: true},
	})

	for pager.More() {
		res, err := pager.NextPage(s.ctx)
		if err != nil {
			log.Printf("Got error [%s]", err)
			continue
		}

		for _, blob := range res.Segment.BlobItems {
			archivesInfo[*blob.Name] = archiveMetadata{
				SizeInBytes: *blob.Properties.ContentLength,
				StoredIn:    []string{Azure.toString()},
			}
		}
	}

	return archivesInfo
}

func (s AzureStore) PutArchive(archiveName string, payload []byte, dest *storageType) error {
	if dest == nil {
		return fmt.Errorf("no destination specified for Azure store")
	}

	if *dest == Azure {
		if archiveName == "" {
			return errors.New("archive name cannot be empty")
		}

		uploadOptions := azblob.UploadBufferOptions{}

		_, err := s.Client.UploadBuffer(s.ctx, s.container, archiveName, payload, &uploadOptions)
		if err != nil {
			return err
		}

		// Dumb workaround because we can't set access tier directly when uploading (xref: https://stackoverflow.com/a/55899242)
		// Have to make a second call to explicitely set the tier
		// Call has to be an HTTP call because sdk doesn't support this operation
		// Thus the need of the extra SAS token for the Azure store setup
		err = s.setAccessTier(archiveName)
		if err != nil {
			return err
		}

		return nil
	}

	log.Printf("wrong destination specified for Azure store [%s], skipping upload", dest.toString())
	return nil
}

func (s AzureStore) DeleteArchive(archiveName string) error {
	_, err := s.Client.DeleteBlob(s.ctx, s.container, archiveName, nil)

	return err
}

func (s *AzureStore) setAccessTier(archiveName string) error {
	endpoint := fmt.Sprintf("%s/%s/%s?comp=tier&%s", s.serviceURL, s.container, archiveName, s.SASToken)

	request, err := http.NewRequest(http.MethodPut, endpoint, nil)
	if err != nil {
		return err
	}
	request.Header.Set("x-ms-date", time.Now().UTC().GoString())
	request.Header.Set("x-ms-access-tier", "Archive")
	request.Header.Set("x-ms-version", "2021-12-02")

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return err
	} else if res.StatusCode >= 400 {
		return fmt.Errorf("got error %s", res.Status)
	}

	return nil
}
