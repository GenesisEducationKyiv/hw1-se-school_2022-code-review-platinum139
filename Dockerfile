FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o main cmd/bitcoin-service/main.go

EXPOSE 8000

CMD [ "./main" ]
