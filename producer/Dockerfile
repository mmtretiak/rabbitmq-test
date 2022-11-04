#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/rabbitmqtest-producer
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/rabbitmqtest-producer -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/rabbitmqtest-producer /rabbitmqtest-producer
ENTRYPOINT /rabbitmqtest-producer
LABEL Name=rabbitmqtest-producer Version=0.0.1
