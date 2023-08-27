package main

import (
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("mqtt_client")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	topic := "file_topic"
	filePath := "arquivo.txt"

	file, err := os.Open(filePath)
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

	token := client.Publish(topic, 0, false, fileBytes)
	token.Wait()

	client.Disconnect(500)
	fmt.Println("File sent:", filePath)
}
