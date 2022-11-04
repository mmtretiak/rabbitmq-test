#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/rabbitmqtest-consumer
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/rabbitmqtest-consumer -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/rabbitmqtest-consumer /rabbitmqtest-consumer
ENTRYPOINT /rabbitmqtest-consumer
LABEL Name=rabbitmqtest-consumer Version=0.0.1
