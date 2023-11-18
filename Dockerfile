FROM docker.io/golang:1.21

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/archiver ./cmd/api

EXPOSE 8080

CMD ["/app/archiver"]