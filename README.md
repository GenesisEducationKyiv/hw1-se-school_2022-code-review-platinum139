# Bitcoin Service

This application provides an API for sending bitcoin rate to subscribers.

## API

1. GET /rate

Returns current BTC-to-UAH rate.

2. POST /subscribe

Subscribes provided email to mailing current rate.

3. POST /sendEmails

Sends email containing current rate to all subscribers.

## How to run service

Make sure you have Docker installed.

1. Clone the repository to your computer.
2. Move to the project folder.
3. Run following command to build docker image
```
docker build --tag bitcoin-svc .
```
4. Run following command to run docker container.
```
docker run --name bitcoin-service -p 8000:8000 bitcoin-svc
```
