# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.22 as builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app statically
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Run stage
FROM scratch

WORKDIR /

# Copy the binary from the builder stage
COPY --from=builder /app/main /main

EXPOSE 8080

CMD ["/main"]
