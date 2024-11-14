FROM golang:1.23-alpine as builder

WORKDIR /app

COPY ./src  /app

RUN go mod tidy
RUN go build -o ./dway .

EXPOSE 8080

CMD ["./dway", "--port=8080", "--config=/etc/dway/config.json"]  