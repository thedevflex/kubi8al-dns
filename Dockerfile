FROM golang:1.24.1 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN mkdir -p /app/.cache
ENV GOCACHE=/app/.cache GOTMPDIR=/app/.cache
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kubi8al-dns ./cmd/...

# Final image
FROM alpine:3.19
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/kubi8al-dns /usr/local/bin/kubi8al-dns

ENTRYPOINT ["/usr/local/bin/kubi8al-dns"]

EXPOSE 8080