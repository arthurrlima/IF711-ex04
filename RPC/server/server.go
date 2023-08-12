package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
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

type FileTransferService struct{}

func (s *FileTransferService) ProcessRequestBytes(chunk FileChunk, reply *bool) error {

	fmt.Println("Recebendo arquivo " + chunk.FileInfo.Filename + " de tamanho " + chunk.FileInfo.Size + " bytes")

	timestamp := time.Now().Format("20060102150405.000000")
	uniqueFileName := timestamp + "_" + chunk.FileInfo.Filename

	outputPath := "files/" + uniqueFileName

	outputFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = outputFile.WriteAt(chunk.Data, chunk.Offset)
	if err != nil {
		return err
	}

	fmt.Println("Enviando dados do arquivo!")

	*reply = true
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
