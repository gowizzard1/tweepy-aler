# Use an official Golang image as the base image
FROM golang:latest

# Set the working directory in the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go code
RUN go build -o tweepy-alert .

# Set the command to run when the container starts
CMD ["./tweepy-alert"]
