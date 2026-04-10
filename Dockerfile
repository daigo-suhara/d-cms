# Stage 1: Build
FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/migrate ./cmd/migrate

# Stage 2: Runtime
FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/bin/server ./server
COPY --from=builder /app/bin/migrate ./migrate
COPY --from=builder /app/templates ./templates

EXPOSE 8080

CMD ["./server"]
