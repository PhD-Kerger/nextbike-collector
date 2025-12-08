# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ./internal ./internal
COPY ./main.go ./main.go
COPY go.mod go.sum ./

RUN go mod download

RUN go build -ldflags="-s -w -buildid=" -trimpath -o nextbike-collector .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy only the binary from the builder stage
COPY --from=builder /app/nextbike-collector /app/nextbike-collector

CMD ["/app/nextbike-collector", "values.yaml"]

