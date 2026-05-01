FROM golang:1.25.0-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go test -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/go-api-server ./cmd/main.go

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S app -G app

COPY --from=builder /bin/go-api-server /usr/local/bin/go-api-server

USER app

EXPOSE 8080

CMD ["go-api-server"]
