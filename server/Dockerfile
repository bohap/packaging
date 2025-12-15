FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first (to cache dependencies)
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary.
# We disable CGO to ensure a static binary that runs on Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port
EXPOSE 8080

# Run the binary
CMD ["./main"]