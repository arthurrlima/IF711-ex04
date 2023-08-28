package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Data interface {
	Bytes() []byte
	Content() string
}

func main() {
	brokerAddress := "localhost:1883" // Change this to your broker's address
	topic := "file_upload_topic"      // Change this to the desired topic

	// Read the text file
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

	// Set up MQTT client options
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerAddress)
	client := mqtt.NewClient(opts)

	// Connect to the broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Error connecting to MQTT broker:", token.Error())
		return
	}
	defer client.Disconnect(250)

	for i := 0; i < 10000; i++ {
		start := time.Now()
		// Publish the file content to the topic
		token := client.Publish(topic, 0, false, fileBytes)
		token.Wait()
		end := time.Since(start)

		record := []string{strconv.FormatInt(end.Milliseconds(), 10)}

		f, _ := os.OpenFile("runlog.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		defer f.Close()

		w := csv.NewWriter(f)
		defer w.Flush()
		w.Write(record)
	}

	// Publish the file content to the topic

	fmt.Println("File sent to server via MQTT.")
}
