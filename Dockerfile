# 1) Сборка
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o todo-server main.go

# 2) Минимальный рантайм
FROM alpine:latest
WORKDIR /app
# нужно CA для HTTPS, но у нас небходимо?
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/todo-server .
COPY web ./web
EXPOSE 7540
ENV TODO_PORT=7540
# Если нужно, можно задать:
# ENV TODO_DBFILE=/data/scheduler.db
# ENV TODO_PASSWORD=secret
CMD ["./todo-server"]