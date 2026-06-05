package client

import (
	"bufio"
	"fmt"
	"net"

	// "os"
	"strconv"
	"strings"
)

func VideoExists() {
	// verifica se o id do trailer existe
	// se conectando com o servidor
}

func DownloadTCP(object string, request_type int, video_id int) ([]byte, error) {

	fmt.Println("\n\n" + object + "\n\n")

	conn, err := net.Dial("tcp", "192.168.15.7:8080")
	if err != nil {
		fmt.Println("Erro ao se conectar no servidor: ", err)
		return nil, err
	}

	// if(request_type == GET_SEGMENT){
	// 	_, err = conn.Write([]byte(strconv.Itoa(GET_SEGMENT) + "\n"))
	// 	_, err = conn.Write([]byte(strconv.Itoa(video_id) + "\n"))
	// 	_, err = conn.Write([]byte(object + "\n")) // Envia a mensagem de um segmentos desejado para o servidor
	// }

	if request_type == LIST_VIDEOS {
		_, err = conn.Write([]byte(strconv.Itoa(LIST_VIDEOS) + "\n"))
		fmt.Println(LIST_VIDEOS)
	}

	if request_type == GET_SEGMENT {
		_, err = conn.Write([]byte(strconv.Itoa(GET_SEGMENT) + "\n"))
		_, err = conn.Write([]byte(strconv.Itoa(video_id) + "\n"))
		_, err = conn.Write([]byte(object + "\n")) // Envia a mensagem de um segmentos desejado para o servidor
	}

	if request_type == GET_MANIFEST {
		_, err = conn.Write([]byte(strconv.Itoa(GET_MANIFEST) + "\n"))
		_, err = conn.Write([]byte(strconv.Itoa(video_id) + "\n"))
	}
	if request_type == GET_THUMBNAIL {
		_, err = conn.Write([]byte(strconv.Itoa(GET_THUMBNAIL) + "\n"))
		_, err = conn.Write([]byte(strconv.Itoa(video_id) + "\n"))
	}
	fmt.Println("oiii")
	connReader := bufio.NewReader(conn) // Reader da conexao
	fmt.Println("antesdoerro")

	// objeto desejado para o servidor
	if err != nil {
		fmt.Println("\tFalha ao enviar nome do objeto desejado")
		conn.Close()
		return nil, err
	}

	msg, _ := connReader.ReadString('\n') // lê READY ou ERROR
	fmt.Println("opaaa" + msg)
	msg = strings.TrimSpace(msg)
	if msg == "ERROR" {
		fmt.Println("\tObjeto nao existe")
		conn.Close()
		return nil, err
	}

	sizeStr, _ := connReader.ReadString('\n')
	sizeStr = strings.TrimSpace(sizeStr)
	size, _ := strconv.Atoi(sizeStr)

	buffer := make([]byte, size)

	var received int

	for received < size {
		n, err := connReader.Read(buffer[received:])
		if err != nil {
			return nil, err
		}
		received += n
	}

	conn.Close()
	return buffer, nil
}
