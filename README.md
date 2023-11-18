# Archiver

## Description

An API to archive and store files.

## Deployment

You can build and run the application with Docker:

```bash
# Build the image
docker build -t archiver .

# Run the container
docker run -d -p 8080:8080 archiver
```

You also have an example of a Kubernetes deployment in the `k8s` folder. Make sure to change the values of the environment variables in the Secret.
