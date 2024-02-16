package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	connectRabbitMQ, conectionOpenError := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(conectionOpenError, "Failed to connect to RabbitMQ")

	defer connectRabbitMQ.Close()
	fmt.Println("Connected to RabbitMQ")

	channelRabbitMQ, channelConnectError := connectRabbitMQ.Channel()
	FailOnError(channelConnectError, "Failed to connect to RabbitMQ")
	defer channelRabbitMQ.Close()

	messages, consumeMessageError := channelRabbitMQ.Consume(
		"first_api", // queue name
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no local
		false,       // no wait
		nil,         // arguments
	)
	FailOnError(consumeMessageError, "Failed to connect to RabbitMQ")

	log.Println("Waiting for messages")

	// Make a channel to receive messages into infinite loop.
	forever := make(chan bool)

	go func() {
		for message := range messages {
			// For example, show received message in a console.
			log.Printf(" > Received message: %s\n", message.Body)
		}
	}()
	// waits for ever since we are not sending any data
	<-forever
}
