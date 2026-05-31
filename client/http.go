package client

import (
	"errors"
	"fmt"
	"net/http"
)

func StartHTTPServer() {
	//servidor HTTP local 

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./client/web")))

	mux.HandleFunc("/video", VideoHandler)
	// mux.HandleFunc("/video/{id}", VideoHandler_v2)
	
	srv := &http.Server{
		Addr: "localhost:3000",
		Handler: mux,
	}

	fmt.Println("Servidor HTTP iniciado em http://localhost:3000/video")
	
	err := srv.ListenAndServe();
	if err != nil{	
		if !errors.Is(err, http.ErrServerClosed){
			panic(err)
		}
	}

}

// func VideoHandler_v2(w http.ResponseWriter, r *http.Request){
// 	id := r.PathValue("id")
	
// 	if id == "" {
// 		http.Error(w, "Nenhum vídeo carregado", http.StatusNotFound)
// 		return
// 	}

// 	// precisa busccar o video,etc
// }

func VideoHandler(w http.ResponseWriter, r *http.Request) {
	// Pensei em seguir esse código que pelo url da página
	// ve o id do video e retorna ele
	// mas precisa implementar
	// por enquanto deixei baixando um vídeo fixo!!!
	// mas a inteção é implementar esse get da url

	


	// http.Error(w, "Nenhum vídeo carregado", http.StatusNotFound)

	FetchVideo("trailer2.mp4\n")

	//define o tipo do conteúdo
	w.Header().Set("Content-Type", "video/mp4")

	// envia o arquivo para o navegador
	http.ServeFile(w, r, "saida/trailer2.mp4") // QUANDO A FUNÇÃO FETCH RETORNA, ja vai ter baixado o video
}

// exec.Command("cmd", "/c", "start", "http://localhost:3000/video").Start()  TEM que colocar de volta mas não sei onde
