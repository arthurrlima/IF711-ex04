package main

import (
	"fmt"
	"net"
	"net/rpc"
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

type FileData struct {
	Filename string
	Size     string
	Origin   int
}
type FileChunk struct {
	FileInfo FileData
	Offset   int64
	Data     []byte
}

type FileTransferService struct{}

var fileCounter = 0

func (s *FileTransferService) ProcessRequestBytes(chunk FileChunk, reply *string) error {
	fileCounter++

	fmt.Println("Recebendo arquivo " + chunk.FileInfo.Filename + " de tamanho " + chunk.FileInfo.Size + " bytes")

	timestamp := time.Now().Format("20060102150405.000000")
	uniqueFileName := timestamp + "_" + strconv.FormatInt(int64(chunk.FileInfo.Origin), 10) + "_" + chunk.FileInfo.Filename

	newFile, err := os.Create("files/" + sanitizeFileName(uniqueFileName))
	if err != nil {
		panic(err)
	}

	fmt.Println("Enviando dados do arquivo!")
	defer newFile.Close()

	_, err = newFile.WriteAt(chunk.Data, chunk.Offset)
	if err != nil {
		return err
	}

	fmt.Println("Enviando dados do arquivo!")

	*reply = "O arquivo foi recebido com sucesso!"
	return nil
}

func main() {
	fmt.Println("Servidor em execução...")

	fileTransferService := new(FileTransferService)
	rpc.Register(fileTransferService)

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
		go rpc.ServeConn(conn)
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
