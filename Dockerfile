# Step 1: Build the binary with the latest Go version
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Copy the project files into the container
COPY . .

# Build the Go binary
RUN go build -o playpi main.go

# Step 2: Create the final image
FROM alpine:latest
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/playpi .

# Expose ports (adjust based on your APIs)
EXPOSE 8080 8081 8082 8083 8084 8085 8086

# Default command to run the binary (can be overridden in docker run)
ENTRYPOINT ["./playpi"]
