# Use the official Golang image as the base image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Install swagger for documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go get -u github.com/swaggo/swag/cmd/swag

# Swag init to generate swagger docs on the working directory
RUN cd /app && swag init


# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080


# Command to run the executable
CMD ["./main"]