// socket-client project main.go
package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	ServerHost = "localhost"
	ServerPort = "1313"
	BUFFERSIZE = 1024
)

func main() {
	rep := make([]byte, 1024)
	// retorna o endereço do endpoint UDP
	addr, err := net.ResolveUDPAddr("udp", ServerHost+":"+ServerPort)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// conecta ao servidor -- não cria uma conexão
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// envia dado
	time.Sleep(21)
	sendFileToServer(conn)

	// recebe resposta do servidor
	_, test, err := conn.ReadFromUDP(rep)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmt.Println(test, " -> ", string(rep))
	// fecha conexão
	defer conn.Close()
}

func sendFileToServer(conn *net.UDPConn) {
	file, err := os.Open("arquivo.txt")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)

	fmt.Println("Enviando nome e tamanho do arquivo!")
	fmt.Println(fileName + ", " + fileSize)

	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Início do upload do arquivo!")
	var bufferString string

	for {
		n, err := file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		fmt.Println(bytes.NewBuffer(sendBuffer[:n]).String())
		bufferString = bytes.NewBuffer(sendBuffer[:n]).String()

	}

	_, err = conn.Write([]byte(fileSize + fileName + bufferString))
	if err != nil {
		fmt.Println("Erro no envio do tamanho e nome do arquivo para o servidor:", err.Error())
		return
	}

	fmt.Println("O upload do arquivo foi concluído! Fechando conexão..." + "filesize: " + fileSize)
}

func fillString(returnString string, toLength int) string {
	for {
		lengthString := len(returnString)
		if lengthString < toLength {
			returnString = returnString + ":"
			continue
		}
		break
	}
	return returnString
}
