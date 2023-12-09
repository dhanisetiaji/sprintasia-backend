# Use the official Golang image as the base image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Expose port 8080 to the outside world
EXPOSE 8080

# # run migrations
# RUN go run gotham -migrate

# # run seeds
# RUN go run gotham -seed

# run the application / migration
CMD ["go", "run", "gotham", "-migration"]