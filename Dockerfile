FROM golang:1.23.2 AS build

WORKDIR /app
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/app

FROM alpine:3.13

COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ=Asia/Jakarta
ENV ZONEINFO=/zoneinfo.zip

COPY --from=build /app/bin/app /app

EXPOSE 3021
ENTRYPOINT ["/app"]

