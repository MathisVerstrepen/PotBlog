# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
COPY . .

RUN templ generate

RUN go mod tidy
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy static assets
COPY --from=builder /app/assets ./assets

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]