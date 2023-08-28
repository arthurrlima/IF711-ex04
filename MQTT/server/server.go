package main

import (
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func saveToFile(filename string, content []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	return err
}

func onMessageReceived(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	fileContent := msg.Payload()

	timestamp := time.Now().Format("20060102150405.000000")
	// Save the received content to a file
	err := saveToFile("files/"+timestamp+"_"+"received_file.txt", fileContent)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}

	fmt.Println("File saved:", topic)
}

func main() {
	brokerAddress := "localhost:1883" // Change this to your broker's address
	topic := "file_upload_topic"      // Change this to the desired topic

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

	// Subscribe to the topic
	go client.Subscribe(topic, 0, onMessageReceived)

	fmt.Println("Server is listening for incoming files via MQTT.")
	select {} // Keep the server running
}
