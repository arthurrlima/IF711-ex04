package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/streadway/amqp"
)

func sendFile(conn *amqp.Connection, channel *amqp.Channel, clientID int) {

	file, err := os.Open("arquivo.txt")
	if err != nil {
		log.Fatal("Error opening the file:", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal("Error getting file info:", err)
	}

	fileBytes := make([]byte, fileInfo.Size())
	_, err = file.Read(fileBytes)
	if err != nil {
		log.Fatal("Error reading the file:", err)
	}
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

	// Publish the file content to the queue
	err = channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        fileBytes,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File sent to server")
}

func main() {

	var wg sync.WaitGroup
	numClients := 10000

	// Create a wait group to synchronize goroutines
	wg.Add(numClients)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	for i := 0; i < numClients; i++ {
		go func(clientID int) {
			defer wg.Done()
			sendFile(conn, ch, clientID)
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}
