// socket-client project main.go
package main

import (
	"bufio"
	"encoding/csv"
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

	// estabelece conexão
	conn, err := net.Dial("tcp", ServerHost+":"+ServerPort)
	if err != nil {
		panic(err)
	}

	for n := 0; n < 1000; n++ {
		start := time.Now()

		// envia dado
		sendFileToClient(conn, value)

		// recebe resposta do servidor
		rep, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		fmt.Print(rep)

		end := time.Since(start)
		record := []string{strconv.FormatInt(end.Milliseconds(), 10)}

		f, err := os.OpenFile("runlog.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		defer f.Close()

		w := csv.NewWriter(f)
		defer w.Flush()

		w.Write(record)

	}

	// fecha conexão
	// defer conn.Close()
}

func sendFileToClient(conn net.Conn, value int) {

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
	fileOrigin := fillString(strconv.Itoa(value), 64)

	fmt.Println("Enviando nome e tamanho do arquivo!")

	_, err = conn.Write([]byte(fileSize))

	if err != nil {
		fmt.Println("Erro no envio do tamanho do arquivo para o servidor:", err.Error())
	}

	_, err = conn.Write([]byte(fileName))
	if err != nil {
		fmt.Println("Erro no envio do nome do arquivo para o servidor:", err.Error())
	}

	_, err = conn.Write([]byte(fileOrigin))
	if err != nil {
		fmt.Println("Erro no envio do numero do cliente para o servidor:", err.Error())
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
