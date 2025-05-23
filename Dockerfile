
FROM golang:1.24.2-alpine

RUN apk update && apk upgrade --no-cache

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod  ./
RUN go mod download

COPY . .

CMD ["go", "test", "./...", "-cover"]
