package main

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func receiveFile(conn *amqp.Connection, channel *amqp.Channel) {
	// Create a queue to receive files.
	q, err := channel.QueueDeclare("myqueue", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Consume messages from the queue.
	messages, err := channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for message := range messages {
		// Save the file to a directory.
		filename := message.Body
		file, err := os.Create("files/" + string(filename))
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

	fmt.Println(ch)
	receiveFile(conn, ch)

}
