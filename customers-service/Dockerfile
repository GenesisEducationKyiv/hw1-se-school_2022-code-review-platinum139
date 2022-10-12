FROM golang:1.18-alpine as builder

WORKDIR /go/src/customers-service

COPY go.mod .
COPY go.sum .

RUN go mod download && go mod verify

COPY . .

RUN go build -o ./app ./cmd/customers-service \
    && chmod +x ./app

FROM alpine
WORKDIR /
COPY --from=builder /go/src/customers-service/app .
COPY --from=builder /go/src/customers-service/.env .

ENTRYPOINT [ "./app" ]
