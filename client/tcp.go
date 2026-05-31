package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func VideoExists(){
	// verifica se o id do trailer existe
	// se conectando com o servidor
}

func FetchVideo(video string) {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao se conectar no servidor: ", err)
		return
	}

	connReader := bufio.NewReader(conn) // Reader da conexao
	_, err = conn.Write([]byte(video)) // Envia a mensagem de um video desejado para o servidor
	// video desejado para o servidor
	if err != nil {
		fmt.Println("\tFalha ao enviar nome do vídeo desejado")
		return
	}

	_, err = conn.Write([]byte(video)) // Envia a mensagem de um 

	msg, _ := connReader.ReadString('\n') // lê READY ou ERROR
	msg = strings.TrimSpace(msg)

	if msg == "ERROR" {
		fmt.Println("\tVideo nao existe")
		return
	}

	sizeStr, _ := connReader.ReadString('\n')
	sizeStr = strings.TrimSpace(sizeStr)
	size, _ := strconv.Atoi(sizeStr)

	video = strings.TrimSpace(video)
	path := "saida/" + video

	out, err := os.Create(path)
	if err != nil {
		fmt.Println("\tErro ao criar o arquivo localmente: ", err)
		return
	}


	fmt.Println("Baixando/streaming do vídeo...")

	//abre navegador automaticamente
	
	buffer := make([]byte, 4096)

	var received int

	for received < size {
		n, _ := connReader.Read(buffer)
		out.Write(buffer[:n])
		received += n
	}

	out.Close()
	fmt.Println("\tDownload finalizado!")
	conn.Close()

}

