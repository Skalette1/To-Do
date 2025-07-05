FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/task-planner ./cmd/server

FROM ubuntu:latest
COPY --from=builder /app/task-planner /task-planner
COPY web /web
CMD ["/task-planner"]
