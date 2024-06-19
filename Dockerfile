# Use the official Golang image as a base
FROM golang:1.22.4

# Set the working directory
WORKDIR /app

# Copy the Go application files
COPY . .

# Download the Go module dependencies
RUN go mod tidy

# Build the Go application
RUN go build -o main .

# Set the entry point for the container
CMD ["./main"]
