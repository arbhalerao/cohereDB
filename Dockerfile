FROM golang:1.23

WORKDIR /cohereDB

COPY . .

RUN go mod tidy

RUN go build -o cohereDB ./cmd/main.go

EXPOSE 8080

ENTRYPOINT ./cohereDB -config=$CONFIG_FILE -container
