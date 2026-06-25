FROM golang:1.26-alpine3.24 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.1/migrate.linux-amd64.tar.gz | tar xvz
        

FROM alpine:3.24
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY --from=builder /app/db/migration ./migration
COPY --from=builder /app/start.sh ./start.sh
COPY --from=builder /app/wait-for.sh ./wait-for.sh

EXPOSE 8080
CMD ["./main"]
ENTRYPOINT ["/app/start.sh"]