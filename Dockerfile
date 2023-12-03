# Use the official Go image as the base image
FROM golang:1.16.3-alpine3.13

# Set the working directory inside the container
# create temp dir for build app
WORKDIR /app

COPY . .


RUN go mod download
RUN go build -o main ./cmd/main.go

EXPOSE 8080

# Command to run the Go application
CMD ["/app/main"]
