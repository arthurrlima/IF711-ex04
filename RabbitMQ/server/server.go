package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

func receiveFile(content amqp.Delivery, count int) {

	clientId := content.MessageId

	timestamp := time.Now().Format("20060102150405.000000")
	uniqueFileName := timestamp + "_" + clientId + "_" + "arquivo.txt"

	file, err := os.Create("files/" + uniqueFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	file.Write(content.Body)

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

	q, err := ch.QueueDeclare(
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
	messages, err := ch.Consume(
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
		fmt.Println("Waiting for messages...")

		go receiveFile(message, count)

	}

}
