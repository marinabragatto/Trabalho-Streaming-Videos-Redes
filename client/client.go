package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"strconv"
)

func main(){
	// Conecta um cliente em yn servidor local na porta 8080
	conn, err :=  net.Dial("tcp", "localhost:8080")
	if err != nil{
		fmt.Println("Falha ao se conectar no servidor: ", err)
		return
	}
	defer conn.Close()

	stdinReader := bufio.NewReader(os.Stdin) // Reader do terminal
	connReader := bufio.NewReader(conn) // Reader da conexao

	
	os.MkdirAll("saida", 0755) // Se não existir o pasta saída, cria

	for{
		fmt.Print("Digite o video: ")
		video, _ := stdinReader.ReadString('\n') // Le a mensagem (do terminal) (considera que ela já vem com o \n que o servidor precisa)
		
		_ , err := conn.Write([]byte(video)) // Envia a mensagem de um video desejado para o servidor
		if err != nil{
			fmt.Println("\tFalha ao enviar nome do vídeo desejado")
			return
		}
		
		msg, _ := connReader.ReadString('\n') // lê READY ou ERROR
		msg = strings.TrimSpace(msg)

		if msg == "ERROR"{
			fmt.Println("\tVideo nao existe")
			return
		}

		sizeStr, _ := connReader.ReadString('\n')
		sizeStr = strings.TrimSpace(sizeStr)
		size, _ := strconv.Atoi(sizeStr)
		
		video = strings.TrimSpace(video)
		path := "saida/" + video
		
		out, err := os.Create(path)
		if err != nil{
			fmt.Println("\tErro ao criar o arquivo localmente: ", err)
			return
		}
		

		buffer := make([]byte, 4096)

		var received int 

		for received < size{
			n, _ := connReader.Read(buffer)
			out.Write(buffer[:n])
			received += n
		}

		out.Close()
		fmt.Println("\tDownload finalizado!")
	}


}