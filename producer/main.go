package main

import (
	"fmt"
	"os"
	"strconv"
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

	replicasCount, err := strconv.Atoi(os.Getenv("REPLICAS_COUNT"))
	if err != nil {
		replicasCount = 10
	}

	rabbitmqServer := "localhost"
	rabbitmqServerStr := os.Getenv("RABBITMQ_SERVER")
	if rabbitmqServerStr != "" {
		rabbitmqServer = rabbitmqServerStr
	}

	msgRate := 1000
	msgRateStr := os.Getenv("MSG_RATE")
	if msgRateStr != "" {
		msgRate64, _ := strconv.ParseInt(msgRateStr, 0, 64)
		msgRate = int(msgRate64)
	}

	msgCount := 100
	msgBurstStr := os.Getenv("MSG_COUNT")
	if msgBurstStr != "" {
		msgBurst64, _ := strconv.ParseInt(msgBurstStr, 0, 64)
		msgCount = int(msgBurst64)
	}

	sleepMs := 1000 / (msgRate / msgCount)

	var conn *amqp.Connection
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

	for i := 0; i < replicasCount; i++ {
		go producerForReplica(ch, hubCount, msgCount, i, sleepMs)
	}
}

func producerForReplica(ch *amqp.Channel, hubCount, msgCount, replica, sleepMs int) {
	for i := 0; i < msgCount; i++ {
		sendMsgToAllHubs(ch, hubCount, replica, sleepMs)
	}
}

func sendMsgToAllHubs(ch *amqp.Channel, hubCount, replica, sleepMs int) {
	for i := 0; i < hubCount; i++ {
		var q amqp.Queue
		var err error

		for {
			q, err = ch.QueueDeclare(
				fmt.Sprintf("%v-hubNum-%v", replica, i),
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

		body := "Hello World!"
		err = ch.Publish("",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})

		time.Sleep(time.Millisecond * time.Duration(sleepMs))
	}
}
