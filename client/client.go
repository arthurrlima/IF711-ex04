// socket-client project main.go
package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

const (
	ServerHost = "localhost"
	ServerPort = "1313"
	ServerType = "tcp"
	BUFFERSIZE = 1024
)

func main() {

	// estabelece conexão
	conn, err := net.Dial(ServerType, ServerHost+":"+ServerPort)
	if err != nil {
		panic(err)
	}

	// envia dado
	sendFileToClient(conn)

	// recebe resposta do servidor
	rep, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmt.Print(rep)

	// fecha conexão
	defer conn.Close()

}

func sendFileToClient(conn net.Conn) {

	file, err := os.Open("arquivo.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)

	fmt.Println("Enviando nome e tamanho do arquivo!")

	_, err = conn.Write([]byte(fileSize))
	if err != nil {
		fmt.Println("Erro no envio do tamanho do arquivo para o servidor:", err.Error())
	}

	_, err = conn.Write([]byte(fileName))
	if err != nil {
		fmt.Println("Erro no envio do nome do arquivo para o servidor:", err.Error())
	}

	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Inicio do upload do arquivo!")

	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		conn.Write(sendBuffer)
	}
	fmt.Println("O upload do arquivo foi concluido! Fechando conexão...")

}

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}
