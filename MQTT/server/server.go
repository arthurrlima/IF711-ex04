package main

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func receiveFile(conn *amqp.Connection, channel *amqp.Channel) {
	// Create a queue to receive files.
	q, err := channel.QueueDeclare(
		"file_queue", // queue name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	// Consume messages from the queue
	messages, err := channel.Consume(
		q.Name, // queue name
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err)
	}

	count := 0

	for message := range messages {
		// Save the file to a directory.
		filename := message.Body
		count++
		file, err := os.Create("files/" + string(filename) + "_" + fmt.Sprintf("%d", count) + ".txt")
		if err != nil {
			fmt.Println(err)
			return
		}

		defer file.Close()

		file.Write(message.Body)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Error connecting to RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Error opening channel:", err)
	}
	defer ch.Close()

	fmt.Println("Waiting for messages...")

	receiveFile(conn, ch)

}
