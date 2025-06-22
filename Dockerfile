FROM golang:1.24.1 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kubi8al-dns ./cmd/main.go

# Final image
FROM alpine:3.19
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/kubi8al-dns /usr/local/bin/kubi8al-dns

ENTRYPOINT ["/usr/local/bin/kubi8al-dns"]