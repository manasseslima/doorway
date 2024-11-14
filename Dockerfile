FROM golang:1.16-alpine as builder

WORKDIR /app

COPY src/* .
COPY src/go.mod .

RUN go mod tidy
RUN go build -o ./dway .

EXPOSE 8080

CMD ["./dway", "--port=8080", "--config=/etc/dway/config.json"]  