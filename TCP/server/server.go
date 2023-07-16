package main

import (
	"fmt"
	"io"
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
	fmt.Println("Servidor em execução...")

	server, err := net.Listen("tcp", ServerHost+":"+ServerPort)
	if err != nil {
		fmt.Println("Erro na escuta por conexões:", err.Error())
		os.Exit(1)
	}

	defer server.Close()

	// aguarda conexões
	fmt.Println("Aguardando conexões dos cliente em " + ServerHost + ":" + ServerPort)

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Cliente conectado")

		// cria thread para o cliente
		go processRequestBytes(conn)
	}
}

func processRequestBytes(conn net.Conn) {
	bufferFileSize := make([]byte, 10)
	bufferFileName := make([]byte, 64)

	// Le o tamanho do arquivo
	conn.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	// Le o nome do arquivo
	conn.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	fmt.Println("Recebendo arquivo " + fileName + " de tamanho " + strconv.FormatInt(fileSize, 10) + " bytes")

	timestamp := time.Now().Format("20060102150405.000000")
	uniqueFileName := timestamp + "_" + fileName

	newFile, err := os.Create("files/" + uniqueFileName)
	if err != nil {
		panic(err)
	}

	fmt.Println("Enviando dados do arquivo!")
	defer newFile.Close()
	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, conn, (fileSize - receivedBytes))
			conn.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, conn, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}

	_, err = conn.Write([]byte("Servidor: Arquivo recebido com sucesso!" + "\n"))
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// fecha conexão
	conn.Close()
}
