FROM golang:1.15.2-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
WORKDIR /srv

COPY --from=builder /app/main /srv
COPY --from=builder /app/.env /srv
COPY --from=builder /app/storage/migrations /srv/storage/migrations
EXPOSE 8080
ENTRYPOINT ["./main"]