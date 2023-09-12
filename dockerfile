# Use the official Go image from the DockerHub
FROM golang:1.18

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. 
RUN go mod download

# Copy the source and the .env file from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 3000 to the outside
EXPOSE 3000

# Run the binary program produced by `go install`
CMD ["./main"]
