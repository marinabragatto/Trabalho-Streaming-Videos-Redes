package client

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func StartHTTPServer() {

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

	fmt.Println("Servidor HTTP iniciado em http://localhost:3000/")
	srv := &http.Server{
		Addr:    "localhost:3000",
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
	data, err := DoRequestListVideos()
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

	segment := r.URL.Query().Get("segment")
	id := r.URL.Query().Get("id")
	id_int, _ := strconv.Atoi(id)
	fmt.Println("(" + segment + ")")

	data, err := DoRequestGetSegment(id_int, segment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	w.Write(data)
	fmt.Println("\tSegmento recebido finalizado!")
}

func ManifestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Rota Envio de Manifesto")
	id := r.URL.Query().Get("id")
	id_int, _ := strconv.Atoi(id)

	data, err := DoRequestGetManifest(id_int)

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

	data, err := DoRequestGetThumbnail(id_int)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(data)
}
