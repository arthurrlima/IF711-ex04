// socket-client project main.go
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

const (
	ServerHost = "localhost"
	ServerPort = "1313"
	BUFFERSIZE = 1024
)

type FileData struct {
	Filename string
	Size     string
}

type FileChunk struct {
	FileInfo FileData
	Offset   int64
	Data     []byte
}

func main() {

	for n := 0; n < 1000; n++ {
		start := time.Now()

		// estabelece conexão
		conn, err := rpc.Dial("tcp", ServerHost+":"+ServerPort)
		if err != nil {
			panic(err)
		}

		// envia dado
		sendFileToClient(conn)

		// recebe resposta do servidor
		rep, err := bufio.NewReader(os.Stdin).ReadString('\n')
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

		// fecha conexão
		defer conn.Close()
	}
}

func sendFileToClient(conn *rpc.Client) {

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

	chunkSize := 1024 // Adjust chunk size as needed
	chunkBuffer := make([]byte, chunkSize)

	var offset int64

	fmt.Println("Enviando nome e tamanho do arquivo!")

	fileData := FileData{
		Filename: fileName,
		Size:     fileSize,
	}

	inputFile, err := os.Open("arquivo.txt")
	fmt.Println("Inicio do upload do arquivo!")

	for {
		n, err := inputFile.Read(chunkBuffer)
		if err != nil && err != io.EOF {
			fmt.Println("Erro de Leitura:", err)
		}
		if n == 0 {
			break
		}

		chunk := FileChunk{
			FileInfo: fileData,
			Offset:   offset,
			Data:     chunkBuffer[:n],
		}

		var reply bool
		err = conn.Call("FileTransferService.ProcessRequestBytes", chunk, &reply)
		if err != nil {
			fmt.Println("Erro:", err)
		}

		if !reply {
			fmt.Println("O servidor falhou em fazer o upload do arquivo.")
		}

		offset += int64(n)
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
