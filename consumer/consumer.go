package main

import (
	errors "dabbit/helpers"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Define RabbitMQ server URL.
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	// Create a new RabbitMQ connection.
	connectRabbitMQ, conectionOpenError := amqp.Dial(amqpServerURL)
	errors.FailOnError(conectionOpenError, "Failed to connect to RabbitMQ")
	log.Println("Successfully connected to RabbitMQ")
	defer connectRabbitMQ.Close()

	channelRabbitMQ, channelConnectError := connectRabbitMQ.Channel()
	errors.FailOnError(channelConnectError, "Failed to connect to RabbitMQ")
	defer channelRabbitMQ.Close()

	messages, consumeMessageError := channelRabbitMQ.Consume(
		"QueueService1", // queue name
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no local
		false,           // no wait
		nil,             // arguments
	)
	errors.FailOnError(consumeMessageError, "Failed to connect to RabbitMQ")

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
