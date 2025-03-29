# Stage 1: Build the Go application
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker caching
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the source code
COPY . .

# Build the application
RUN go build -o bot ./cmd/main.go

# Stage 2: Create a lightweight runtime image
FROM debian:bookworm-slim

# Install yt-dlp and dependencies, overriding the externally managed environment
RUN apt-get update && \
    apt-get install -y python3 python3-pip && \
    pip3 install yt-dlp --break-system-packages && \
    rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/bot .

# Copy any additional assets (e.g., Instagram cookies)
COPY ./assets ./assets

# Command to run the bot
CMD ["./bot"]