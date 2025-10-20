FROM golang:1.24.5-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go generate ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN addgroup -S app && adduser -S app -G app

RUN mkdir -p /app/storage
RUN chown -R app:app /app

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Copy environment file - just this time
COPY .env .

USER app

EXPOSE 50051

CMD ["./main"]