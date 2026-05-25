package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

/* Esse código coloca em execução o servidor na porta :8080 do localhost*/
// Seria interessante apresentar um catálogo de filmes (trailers) ou o vídeo que for, para o usuário poder escolher o que ele quer

func main(){
	// Listener : interface objeto que escuta conexões
	listener, err := net.Listen("tcp", ":8080")  // Cria socket tcp
	if err != nil{
		fmt.Println(":Erro ao iniciar o servidor", err) // Trata erro ao criar socket
		return
	}
	defer listener.Close() // Fecha o socket que foi aberto quando a funçao main terminar 

	fmt.Println("Server está escutando na porta 8080")

	// Aceita conexoes em looping
	for{
		conn, err := listener.Accept() // conn : Conexão tcp
		// Accept é bloqueante
		if err != nil {
			fmt.Println("Failed to accept connection: ", err)
			continue
		}

		go handleConnection(conn) 
	}
}

// Responsável por tratar uma conexão
/*
FORMATO DE LEITURA DE CONEXÃO (PROTOCOLO):
  .1	Servidor recebe do cliente o nome do arquivo:
			nome.mp4\n
  .2	Se existe o arquivo:
			READY\n
			SIZE\n (o size é um int)
			BYTES DO ARQUIVO
		Se NÃO existe o arquivo:
			ERROR\n

*/

func handleConnection(conn net.Conn){
	defer conn.Close()
	reader := bufio.NewReader(conn)
	fmt.Print("Nova Conexão\n")
	// Considerando que o cliente envia no formato .mp4 :
	//trailer_michael.mp4\n
	//(o cliente DEVE enviar o \n)
	video, _ := reader.ReadString('\n') // Le a conexão até \n
	video = strings.TrimSpace(video)

	path := "../videos/" + video // Caminho até a pasta de arquivos do servidor
	fmt.Println("\tVídeo Solicitado: " + video) // Imprime o pedido desejado no terminal do servidor 

	file, err := os.Open(path)
	if err != nil {
		fmt.Print("Erro: não existe esse arquivo na base de dados\n")
		conn.Write([]byte("ERROR\n"))
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	
	conn.Write([]byte("READY\n"))
	conn.Write([]byte(fmt.Sprintf("%d\n", size)))
	buffer := make([]byte, 4096)
		
	fmt.Println("\tEnviando o arquivo ...")
	for {
		n, err := file.Read(buffer) // Le do arquivo .mp4 e guarda no buffer
		if n > 0{
			conn.Write((buffer[:n])) // Se não chegou no fim do arquivo de vídeo, continua
		}
		if err != nil{
			break // Se leu tudo (EOF), termina
		}
	}

	fmt.Println("\tArquivo enviado: " , video)

}