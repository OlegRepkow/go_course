FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN GOWORK=off go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY . .
RUN /go/bin/swag init -g cmd/server/main.go
RUN GOWORK=off go build -o server ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
