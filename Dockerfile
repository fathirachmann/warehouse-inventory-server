# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build tools
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o warehouse-server ./

# Run stage
FROM alpine:3.20

WORKDIR /app

# For Fiber / networking (port configurable via env PORT)
EXPOSE 8080

COPY --from=builder /app/warehouse-server /app/warehouse-server

CMD ["/app/warehouse-server"]
