FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/todo-server ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/todo-server /app/todo-server
COPY --from=builder /app/web /app/web

RUN chmod +x /app/todo-server

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db

ENTRYPOINT ["/app/todo-server"]