package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func VideoExists() {
	// verifica se o id do trailer existe
	// se conectando com o servidor
}

func Download(object string) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao se conectar no servidor: ", err)
		return
	}

	connReader := bufio.NewReader(conn)        // Reader da conexao
	_, err = conn.Write([]byte(object + "\n")) // Envia a mensagem de um segmentos desejado para o servidor
	// objeto desejado para o servidor
	if err != nil {
		fmt.Println("\tFalha ao enviar nome do objeto desejado")
		conn.Close()
		return
	}

	msg, _ := connReader.ReadString('\n') // lê READY ou ERROR
	msg = strings.TrimSpace(msg)

	if msg == "ERROR" {
		fmt.Println("\tObjeto nao existe")
		conn.Close()
		return
	}

	sizeStr, _ := connReader.ReadString('\n')
	sizeStr = strings.TrimSpace(sizeStr)
	size, _ := strconv.Atoi(sizeStr)

	path := "./client/segments/" + object

	out, err := os.Create(path)
	if err != nil {
		fmt.Println("\tErro ao criar o arquivo localmente: ", err)
		conn.Close()
		return
	}

	buffer := make([]byte, 4096)

	var received int

	for received < size {
		n, _ := connReader.Read(buffer)
		out.Write(buffer[:n])
		received += n
	}

	out.Close()
	conn.Close()
}

func FetchVideo() {
	os.MkdirAll("./client/segments", os.ModePerm)

	Download("manifest.json")
	fmt.Println("Manifesto recebido com sucesso")

	segments, err := ReadManifest()
	if err != nil {
		fmt.Println("Erro ao ler manifesto", err)
		return
	}

	for _, segment := range segments {
		Download(segment)
		fmt.Println("\tSegmento recebido finalizado!")
	}

	fmt.Println("\tTodos os segmentos foram baixados finalizado!")

}
