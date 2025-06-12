# Start from the official Go image
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go files
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the binary
RUN go build -o app ./cmd/app

# Final image
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/app .

CMD ["/app/app"]
