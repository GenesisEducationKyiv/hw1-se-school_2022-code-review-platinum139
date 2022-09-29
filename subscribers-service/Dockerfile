FROM golang:1.18-alpine as builder

WORKDIR /go/src/subscribers-service

COPY go.mod .
COPY go.sum .

RUN go mod download && go mod verify

COPY . .

RUN go build -o ./app ./cmd/subscribers-service \
    && chmod +x ./app

FROM alpine
WORKDIR /
COPY --from=builder /go/src/subscribers-service/app .
COPY --from=builder /go/src/subscribers-service/.env .

ENTRYPOINT [ "./app" ]
