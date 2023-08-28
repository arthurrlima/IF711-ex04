package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

func sendFile(conn *amqp.Connection, channel *amqp.Channel, clientID int, fileID int) {

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
			MessageId:   fmt.Sprintf("%d", clientID) + "_" + fmt.Sprintf("%d", fileID),
			Body:        fileBytes,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File sent to server")
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./client <value>")
		return
	}

	arg := os.Args[1]
	value, err := strconv.Atoi(arg)
	if err != nil {
		fmt.Println("Invalid argument:", err)
		return
	}

	// Create a wait group to synchronize goroutines

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	for i := 0; i < 10000; i++ {
		start := time.Now()
		sendFile(conn, ch, value, i)
		end := time.Since(start)
		record := []string{strconv.FormatInt(end.Milliseconds(), 10)}

		f, _ := os.OpenFile("runlog.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		defer f.Close()

		w := csv.NewWriter(f)
		defer w.Flush()

		w.Write(record)
	}
}
