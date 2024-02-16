package main

import (
	"context"
	"fmt"
	"log"

	fiber "github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// TODO: connect to rabbitmq
	rabbitConnection, conectionOpenError := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(conectionOpenError, "Failed to connect to RabbitMQ")

	defer rabbitConnection.Close()
	fmt.Println("Connected to RabbitMQ")

	// TODO: open a channel
	channel, channelOpenError := rabbitConnection.Channel()
	FailOnError(channelOpenError, "Failed to open a channel")
	defer channel.Close()
	fmt.Println("Opened a channel")

	// TODO: create a queue
	queue, err := channel.QueueDeclare(
		"first_api",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue")

	// for undeliverable messages that are returned by the server
	channel.NotifyReturn(make(chan amqp.Return))

	//TODO: create API
	app := fiber.New()
	fmt.Println("Successfully instantiated a server")

	app.Get("/", func(c *fiber.Ctx) error {
		message := c.Query("message")

		ctx, cancel := context.WithTimeout(context.Background(), 10000)
		defer cancel()

		messageObj := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		}
		err = channel.PublishWithContext(
			ctx,
			"",
			queue.Name,
			true,
			false,
			messageObj,
		)
		FailOnError(err, "Failed to publish a message")

		fmt.Println("Message sent to: " + queue.Name)
		return c.SendString("Message sent to: " + queue.Name)
	})

	serverStartError := app.Listen(":3000")
	FailOnError(serverStartError, "Failed to listen")

}

//TODO: mock failed messages and handle that
