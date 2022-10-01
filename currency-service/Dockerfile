FROM golang:1.18-alpine as builder

WORKDIR /go/src/currency-service

COPY go.mod .
COPY go.sum .

RUN go mod download && go mod verify

COPY . .

RUN go build -o ./app ./cmd/currency-service \
    && chmod +x ./app

FROM alpine
WORKDIR /
COPY --from=builder /go/src/currency-service/app .
COPY --from=builder /go/src/currency-service/.env .

ENTRYPOINT [ "./app" ]
