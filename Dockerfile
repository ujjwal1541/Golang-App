FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /healthcare-app ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /healthcare-app /app/healthcare-app

# Copy migrations
COPY --from=builder /app/migrations /app/migrations

# Install necessary tools
RUN apk --no-cache add ca-certificates && \
    chmod +x /app/healthcare-app

# Expose port
EXPOSE 8080

# Command to run the application
CMD ["/app/healthcare-app"] 