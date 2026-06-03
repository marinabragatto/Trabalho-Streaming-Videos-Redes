package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

/* Esse código coloca em execução o servidor na porta :8080 do localhost*/
// Seria interessante apresentar um catálogo de filmes (trailers) ou o vídeo que for, para o usuário poder escolher o que ele quer

const (
	LIST_VIDEOS   = 1
	GET_THUMBNAIL = 2
	GET_MANIFEST  = 3
	GET_SEGMENT   = 4
)

func (v Video) get_id() int {
	return v.Id
}

type Video struct {
	Id        int   `json:"id"`
	Nome      string `json:"nome"`
	Thumbnail string  `json:"thumbnail"`
	Manifest  string `json:"manifest"`
}

const num_videos int = 4

func Serve() {

	var videos []Video

	// JSON
	file, err := os.Open("metadata.json")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&videos)
	if err != nil {
		fmt.Println("Erro ao ler JSON", err)
		return
	}
	imprime_videos(videos)

	// Listener : interface objeto que escuta conexões
	listener, err := net.Listen("tcp", "localhost:8080") // Cria socket tcp
	if err != nil {
		fmt.Println(":Erro ao iniciar o servidor", err) // Trata erro ao criar socket
		return
	}
	defer listener.Close() // Fecha o socket que foi aberto quando a funçao main terminar

	fmt.Println("Server está escutando na porta 8080")

	// Aceita conexoes em looping
	for {
		conn, err := listener.Accept() // conn : Conexão tcp
		// Accept é bloqueante
		if err != nil {
			fmt.Println("Failed to accept connection: ", err)
			continue
		}

		go handleConnection(conn, videos)
	}
}

// Responsável por tratar uma conexão
/*
FORMATO DE LEITURA DE CONEXÃO (PROTOCOLO):
  .1	Servidor recebe o tipo de REQUEST
  .1	Servidor recebe do cliente o nome do arquivo:
			nome.mp4\n
  .2	Se existe o arquivo:
			READY\n
			SIZE\n (o size é um int)
			BYTES DO ARQUIVO
		Se NÃO existe o arquivo:
			ERROR\n

*/

func imprime_videos(videos []Video) {

	for i := 0; i < len(videos); i++ {
		fmt.Println(videos[i].Nome)
	}
}

func handleConnection(conn net.Conn, videos []Video) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	fmt.Print("Nova Conexão\n")
	// Considerando que o cliente envia no formato .mp4 :
	//trailer_michael.mp4\n
	//(o cliente DEVE enviar o \n)
	request, _ := reader.ReadString('\n')
	request = strings.TrimSpace(request)
	request_type, _ := strconv.Atoi(request)

	print("SERVIDOR: listando videos")

	if request_type == LIST_VIDEOS {

		print("SERVIDOR: listando videos")
		data, err := json.Marshal(videos)
		if err != nil {
			fmt.Println(err)
			return
		}
		conn.Write([]byte("READY\n"))
		conn.Write([]byte(fmt.Sprintf("%d\n", len(data))))
		conn.Write(data) // deveria segmentar?
		return
	}

	object, _ := reader.ReadString('\n') // Le a conexão até \n
	object = strings.TrimSpace(object)
	fmt.Println(object)
	video_id, _ := reader.ReadString('\n')
	video_id = strings.TrimSpace(video_id)
	video_id_num, _ := strconv.Atoi(video_id)

	if video_id_num > num_videos && video_id_num > 0 { // o video eh de id que nao se conhece
		fmt.Print("Erro: não existe esse arquivo na base de dados\n")
		conn.Write([]byte("ERROR\n"))
		return
	}

	path := ""
	if request_type == GET_MANIFEST {
		path = "./videos/" + video_id + "/manifest.json" // Caminho até a pasta de arquivos do servidor
	}
	if request_type == GET_THUMBNAIL {
		path = "./thumbnail/" + video_id + ".jpg" // Caminho até a pasta de arquivos do servidor

	}
	if request_type == GET_SEGMENT {
		nome_seg, _ := reader.ReadString('\n')
		path = "./videos/" + video_id + "/" + nome_seg // Caminho até a pasta de arquivos do servidor

	}

	fmt.Println("\tVídeo Solicitado: " + path) // Imprime o pedido desejado no terminal do servidor

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo desejado")
		conn.Write([]byte("ERROR\n"))
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
		if n > 0 {
			conn.Write((buffer[:n])) // Se não chegou no fim do arquivo de vídeo, continua
		}
		if err != nil {
			break // Se leu tudo (EOF), termina
		}
	}

	fmt.Println("\tArquivo enviado: ", object)

}
