package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	// "net/http/cookiejar"

	// "os"
	"strconv"
	"strings"
)


func DownloadData(conn net.Conn) ([]byte, error){
	connReader := bufio.NewReader(conn) // Reader da conexao
	msg, _ := connReader.ReadString('\n') // lê READY ou ERROR
	fmt.Println("opaaa"+msg)
	msg = strings.TrimSpace(msg)
	
	if msg == "ERROR" {
		err := errors.New("Operação não existe no servidor!")
		fmt.Println(err)
		return nil, err
	}

	sizeStr, _ := connReader.ReadString('\n') // Le o Size dos dados
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



func DoRequestGetManifest(video_id int) ([]byte, error) {
	conn, err := net.Dial("tcp", "192.168.15.7:8080")
	// conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao se conectar no servidor: ", err)
		return nil, err
	}
	_, err = conn.Write([]byte(strconv.Itoa(GET_MANIFEST) + "\n"))
	_, err = conn.Write([]byte(strconv.Itoa(video_id) + "\n"))

	return DownloadData(conn)
}

func DoRequestGetSegment(video_id int, segment string) ([]byte, error) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao se conectar no servidor: ", err)
		return nil, err
	}
	_, err = conn.Write([]byte(strconv.Itoa(GET_SEGMENT) + "\n"))
	_, err = conn.Write([]byte(strconv.Itoa(video_id) + "\n"))
	_, err = conn.Write([]byte(segment + "\n")) // Envia a mensagem de um segmentos desejado para o servidor

	return DownloadData(conn)
}

func DoRequestGetThumbnail(video_id int) ([]byte, error) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao se conectar no servidor: ", err)
		return nil, err
	}
	_, err = conn.Write([]byte(strconv.Itoa(GET_THUMBNAIL) + "\n"))
	_, err = conn.Write([]byte(strconv.Itoa(video_id) + "\n"))
	
	return DownloadData(conn)
}

func DoRequestListVideos() ([]byte, error) {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao se conectar no servidor: ", err)
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(strconv.Itoa(LIST_VIDEOS) + "\n"))
	fmt.Println(LIST_VIDEOS)
	if err != nil {
		fmt.Println("\tFalha ao enviar Operação Desejada")
		conn.Close()
		return nil, err
	}

	return DownloadData(conn)
	
}


