FROM golang:1.18-alpine as builder

WORKDIR /go/src/logs-printer

COPY go.mod .
COPY go.sum .

RUN go mod download && go mod verify

COPY . .

RUN go build -o ./app ./cmd/logs-printer \
    && chmod +x ./app

FROM alpine
WORKDIR /
COPY --from=builder /go/src/logs-printer/app .
COPY --from=builder /go/src/logs-printer/.env .

ENTRYPOINT [ "./app" ]
