FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o /dragon-chat

ENV PORT=8080

EXPOSE $PORT

CMD ["/dragon-chat"]