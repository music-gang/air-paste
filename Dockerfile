# Use an official Golang runtime as a parent image
FROM golang:1.22

# Metadata
LABEL org.opencontainers.image.source = "https://github.com/music-gang/air-paste"

# Set the working directory inside the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Build the Go app
RUN go build -o main ./cmd

# Environment variable
ENV PORT=8080

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]