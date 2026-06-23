FROM golang:1.26-alpine3.24 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.24
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./main"]