# Use the official Golang base image for compilation
FROM golang:latest as builder

# Set the go binary directory
ENV GOBIN /usr/local/bin/

# Install entr and fd (fd-find is already included in golang:latest)
RUN apt-get update && \
	apt-get install -y entr && \
	rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /app

# Install auto Reloader
RUN go install github.com/bokwoon95/wgo@latest

# Copy the Go module files to the container
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy commands
COPY cmd cmd

# Copy internals
COPY internal internal

# Copy packages
COPY pkg pkg

# Build the binaries
RUN CGO_ENABLED=1 GOOS=linux go install ./cmd/...
