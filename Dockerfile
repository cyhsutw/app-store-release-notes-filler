FROM golang:1.16-alpine

RUN apk update \
    && apk add gcc \
    && apk add musl-dev \
    && apk add sqlite \
    && apk add socat

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /server

EXPOSE 8080

CMD [ "/server" ]
