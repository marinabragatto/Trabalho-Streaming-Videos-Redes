package client

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)
var currentVideoPath string

func Connect_to_server()(){
	// Conecta um cliente em yn servidor local na porta 8080
	conn, err :=  net.Dial("tcp", "localhost:8080")
	if err != nil{
		fmt.Println("Falha ao se conectar no servidor: ", err)
		return
	}
	defer conn.Close()

	stdinReader := bufio.NewReader(os.Stdin) // Reader do terminal

	os.MkdirAll("saida", 0755) // Se não existir o pasta saída, cria

	//servidor HTTP local
	http.HandleFunc("/video", func(w http.ResponseWriter, r *http.Request) {

		if currentVideoPath == "" {
			http.Error(w, "Nenhum vídeo carregado", http.StatusNotFound)
			return
		}

		//define o tipo do conteúdo
		w.Header().Set("Content-Type", "video/mp4")

		//envia o arquivo para o navegador
		http.ServeFile(w, r, currentVideoPath)
	})

	go func() {
		fmt.Println("Servidor HTTP iniciado em http://localhost:3000/video")
		http.ListenAndServe(":3000", nil)
	}()

	for {
		// Conecta um cliente em um servidor local na porta 8080
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			fmt.Println("Falha ao se conectar no servidor: ", err)
			return
		}

		connReader := bufio.NewReader(conn) // Reader da conexao

		fmt.Print("Digite o video: ")
		video, _ := stdinReader.ReadString('\n') // Le a mensagem (do terminal) (considera que ela já vem com o \n que o servidor precisa)

		_, err = conn.Write([]byte(video)) // Envia a mensagem de um video desejado para o servidor
		if err != nil {
			fmt.Println("\tFalha ao enviar nome do vídeo desejado")
			return
		}

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

		//atualiza vídeo atual do localhost
		currentVideoPath = path

		fmt.Println("Baixando/streaming do vídeo...")

		//abre navegador automaticamente
		exec.Command("cmd", "/c", "start", "http://localhost:3000/video").Start()

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

}
