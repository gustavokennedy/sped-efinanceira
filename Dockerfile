FROM golang:latest AS builder

WORKDIR /app

# Define a flag para o build
ENV GOFLAGS="-buildvcs=false"

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/codegangsta/gin@latest

# Define volume
VOLUME ["/app"]

# Define a porta de execução
EXPOSE 8080

# Reinicializa o servidor com gin
ENTRYPOINT ["gin", "-i", "run"]