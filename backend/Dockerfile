# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o freenstagram ./cmd/server/main.go

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/freenstagram .
# Create a directory for image uploads
RUN mkdir uploads

EXPOSE 8080

CMD ["./freenstagram"]
