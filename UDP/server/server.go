package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ServerHost = "localhost"
	ServerPort = "1313"
	BUFFERSIZE = 1024
)

func main() {
	msgFromClient := make([]byte, 1024)

	// resolve server address
	addr, err := net.ResolveUDPAddr("udp", ":"+ServerPort)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// listen on udp port
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// close conn
	// close conn
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	fmt.Println("Servidor UDP aguardando requests...")

	for {

		// receive request
		n, addr, err := conn.ReadFromUDP(msgFromClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		go processRequestBytes(conn, msgFromClient[:n], n, addr)
	}
}

func processRequestBytes(conn *net.UDPConn, msgFromClient []byte, n int, addr *net.UDPAddr) {
	bufferFileSize := msgFromClient[:10]
	bufferFileName := msgFromClient[10:74]

	// Le o tamanho do arquivo
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	// Le o nome do arquivo
	fileName := strings.Trim(string(bufferFileName), ":")

	fmt.Println("Recebendo arquivo " + fileName + " de tamanho " + strconv.FormatInt(fileSize, 10) + " bytes")

	timestamp := time.Now().Format("20060102150405.000000")
	uniqueFileName := timestamp + "_" + fileName

	newFile, err := os.Create("files/" + sanitizeFileName(uniqueFileName))
	if err != nil {
		panic(err)
	}

	fmt.Println("Enviando dados do arquivo!")
	defer newFile.Close()
	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			remainingBufferSize := fileSize - receivedBytes
			if remainingBufferSize <= 0 {
				break
			}
			remainingBuffer := make([]byte, remainingBufferSize)
			n, _, err := conn.ReadFromUDP(remainingBuffer)
			if err != nil {
				fmt.Println("Failed to read UDP packet:", err)
				break
			}
			newFile.Write(remainingBuffer[:n])
			receivedBytes += int64(n)
			break
		}
		n, _, err := conn.ReadFromUDP(msgFromClient)
		if err != nil {
			fmt.Println("Failed to read UDP packet:", err)
			break
		}
		newFile.Write(msgFromClient[:n])
		receivedBytes += BUFFERSIZE
	}

	_, err = conn.WriteToUDP([]byte("Servidor: Arquivo recebido com sucesso!\n"), addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func sanitizeFileName(fileName string) string {
	// Remove any invalid characters from the file name
	invalidChars := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		fileName = strings.ReplaceAll(fileName, char, "")
	}

	return fileName
}
