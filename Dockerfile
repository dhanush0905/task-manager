# Use an official Go runtime as the base image
FROM golang:1.21

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire backend code
COPY . .

# Build the Go app
RUN go build -o main .

# Expose the application port
EXPOSE 8080

# Start the application
CMD ["./main"]
