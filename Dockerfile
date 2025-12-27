# Stage 1: Build the Go application
FROM golang:1.25 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
# CGO_ENABLED=0 prevents the use of Cgo, which is needed for cross-compilation
# GOOS=linux specifies the target operating system as Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api .

# Stage 2: Create the final lightweight image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/api .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
# This will be the main process for the container
CMD ["./api"]