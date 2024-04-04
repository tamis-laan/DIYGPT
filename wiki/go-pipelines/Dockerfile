# Use the official Golang base image for compilation
FROM golang:latest as builder

# Install entr and fd (fd-find is already included in golang:latest)
RUN apt-get update && \
	apt-get install -y entr && \
	rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /app

# Set the go binary directory
ENV GOBIN /usr/local/bin/

# Copy the Go module files to the container
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code to the container
COPY . .

# Build the binaries
RUN CGO_ENABLED=1 GOOS=linux go install ./cmd/...

# # Start a new stage from scratch
# FROM alpine:latest as deployment

# # Set the working directory inside the container
# WORKDIR /app

# # Copy the binary from the builder stage
# # COPY --from=builder /usr/local/bin /
# COPY --from=builder /usr/local/bin/ .

# RUN ls /app
