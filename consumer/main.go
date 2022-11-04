package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	hubCount := 100
	hubCountStr := os.Getenv("HUB_COUNT")
	if hubCountStr != "" {
		hubCount64, _ := strconv.ParseInt(hubCountStr, 0, 64)
		hubCount = int(hubCount64)
	}

	replica := 0
	hostname := os.Getenv("POD_NAME")
	fmt.Println(hostname)
	ss := strings.Split(hostname, "consumer-stateful-set-")
	if len(ss) > 1 {
		replica64, _ := strconv.ParseInt(ss[1], 0, 64)
		replica = int(replica64)
	}
	fmt.Println(replica)

	rabbitmqServer := "localhost"
	rabbitmqServerStr := os.Getenv("RABBITMQ_SERVER")
	if rabbitmqServerStr != "" {
		rabbitmqServer = rabbitmqServerStr
	}

	for i := 0; i < hubCount; i++ {
		go consumer(rabbitmqServer, i, replica)
		time.Sleep(10 * time.Millisecond)
	}

	var forever chan struct{}
	<-forever
}

func consumer(rabbitmqServer string, hubNum int, replica int) {
	var conn *amqp.Connection
	var err error
	for {
		conn, err = amqp.Dial(rabbitmqServer)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		break
	}
	defer conn.Close()

	var ch *amqp.Channel
	for {
		ch, err = conn.Channel()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		break
	}
	defer ch.Close()

	var q amqp.Queue
	for {
		q, err = ch.QueueDeclare(
			fmt.Sprintf("%v-hubNum-%v", replica, hubNum),
			false,
			true,
			false,
			false,
			nil,
		)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		break
	}

	var msgs <-chan amqp.Delivery
	for {
		msgs, err = ch.Consume(
			q.Name,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			continue
		}

		break
	}

	for m := range msgs {
		log.Println("[Hub ", hubNum, "] Received js message w/ type", m.Type, "data", string(m.Body))
		m.Ack(false)
	}
}
