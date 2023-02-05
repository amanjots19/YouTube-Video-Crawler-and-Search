# Use an official Golang image as the base image
FROM golang:latest

# Set the working directory
WORKDIR /app

# Copy the source code from the host to the container
COPY ./cmd/yt-fetch /app

# Build the binary for the "yt-fetch" application
RUN go mod init yt-fetch
RUN go mod tidy
RUN go build -o ./cmd/...

# Set the environment variable to tell the application which port to listen on
ENV PORT 8080

# Expose the port that the application is listening on
EXPOSE 8080

# Set the command to run when the container starts
CMD ["./yt-fetch"]