FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY .env .
RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta
EXPOSE 8080
ENTRYPOINT ["./main"]