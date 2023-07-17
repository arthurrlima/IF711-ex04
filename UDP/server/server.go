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
	bufferFileContent := msgFromClient[74:]

	// Le o tamanho do arquivo
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	// Le o nome do arquivo
	fileName := strings.Trim(string(bufferFileName), ":")

	fmt.Println("Recebendo arquivo " + fileName + " de tamanho " + strconv.FormatInt(fileSize, 10) + " bytes")

	// File content
	fmt.Println(bufferFileContent)

	timestamp := time.Now().Format("20060102150405.000000")
	uniqueFileName := string(timestamp[15:]) + fileName

	newFile, err := os.Create("files/" + sanitizeFileName(uniqueFileName))
	if err != nil {
		panic(err)
	}

	fmt.Println("Enviando dados do arquivo!")
	defer newFile.Close()

	newFile.Write(bufferFileContent)

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
