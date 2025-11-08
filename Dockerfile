FROM golang:1.25 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN  go build -o api main.go

EXPOSE 8080
CMD ["./api"]
