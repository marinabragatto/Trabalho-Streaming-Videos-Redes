package client

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var ip string

func ReadIP() string {
	fmt.Print("Digite aqui o IP do seu SERVIDOR TCP: ")

	reader := bufio.NewReader(os.Stdin)

	ip, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Erro ao ler IP:", err)
		return ""
	}
	ip = strings.TrimSpace(ip)
	return ip
}

func StartHTTPServer() {

	ip = ReadIP()
	if ip == "" {
		return
	}

	mux := http.NewServeMux()

	// Rota para o Home, entrega interface VAZIA
	mux.Handle("/", http.FileServer(http.Dir("./client/web")))
	// Rota de envio de catalogo para o front-end
	mux.HandleFunc("/api/catalogo", CatalogoHandler)
	// Rota para o Video Player
	mux.HandleFunc("/video", VideoHandler)
	// Rota para o streaming via Segmentos
	mux.HandleFunc("/stream", StreamHandler)
	// Rota Auxiliar para envio do Manifesto para o front-end
	mux.HandleFunc("/manifest", ManifestHandler)
	// Rota para a Thumbnail
	mux.HandleFunc("/thumbnail", ThumbnailHandler)

	fmt.Printf("Servidor HTTP iniciado em http://%s:3000/\n", ip)
	srv := &http.Server{
		Addr:    "0.0.0.0:3000",
		Handler: mux,
	}

	err := srv.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}

func CatalogoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Rota Catalogo")
	data, err := DoRequestListVideos(ip)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Println("Voltou do DownloadTCP")

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func VideoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Rota Interface Video-Player")
	http.ServeFile(w, r, "./client/web/video.html")
}

func StreamHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Rota Stream!")
	quality := r.URL.Query().Get("quality")
	quality_int, _ := strconv.Atoi(quality)
	segment := r.URL.Query().Get("segment")
	fmt.Println("SEGMENTO + ", segment)
	id := r.URL.Query().Get("id")
	id_int, _ := strconv.Atoi(id)
	fmt.Println("(" + segment + ")")

	data, err := DoRequestGetSegment(id_int, quality_int, segment, ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if strings.HasSuffix(segment, ".m4s") || strings.HasSuffix(segment, ".mp4") {
		w.Header().Set("Content-Type", "video/mp4")
	}
	w.Write(data)
	fmt.Println("\tSegmento recebido finalizado!")
}

func ManifestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Rota Envio de Manifesto")
	id := r.URL.Query().Get("id")
	id_int, _ := strconv.Atoi(id)

	data, err := DoRequestGetManifest(id_int, ip)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func ThumbnailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Rota Thumbnail")
	id := r.URL.Query().Get("id")
	id_int, err := strconv.Atoi(id)

	data, err := DoRequestGetThumbnail(id_int, ip)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(data)
}
