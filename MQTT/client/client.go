package main

import (
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

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

	// Publish the file content to the topic
	token := client.Publish(topic, 0, false, fileBytes)
	token.Wait()

	fmt.Println("File sent to server via MQTT.")
}
