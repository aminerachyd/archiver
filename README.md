# Archiver

## Description

An API to archive and store files.  
Currently supports Azure Blob Storage and file system storage.

## Deployment

### Docker

You can build and run the application with Docker:

```bash
# Build the image
docker build -t archiver .

# Run the container with Docker
docker run -d -p 8080:8080 \
 -e AZURE_SERVICE_URL=<BLOB_STORE_URL> \ # Url looks like https://<storageaccount>.blob.core.windows.net
 -e AZURE_SAS_TOKEN=<SAS_TOKEN> \
 -e AZURE_TENANT_ID=<TENANT_ID> \
 -e AZURE_CLIENT_ID=<CLIENT_ID> \
 -e AZURE_CLIENT_SECRET=<CLIENT_SECRET> \
 -e AZURE_ARCHIVES_CONTAINER=<YOUR_CONTAINER_IN_STORAGEACCOUNT> \
 -e ARCHIVER_FILESYSTEM_PATH=<LOCAL_PATH> \ # For local storage on disk
archiver

# Check if the API is alive
curl http://localhost:8080/v1/health
```

### Kubernetes

You have an example of a Kubernetes deployment in the `k8s` folder.  
Make sure to change the values of the environment variables in the Secret.
