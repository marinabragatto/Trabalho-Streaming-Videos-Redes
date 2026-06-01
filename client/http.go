package client

import (
	"errors"
	"fmt"
	"net/http"
	"os/exec"
)

func StartHTTPServer() {
	//servidor HTTP local

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./client/web")))

	// Serve os segmentos baixados pelo TCP
	mux.Handle("/segments/", http.StripPrefix("/segments/", http.FileServer(http.Dir("./client/segments"))))

	//THREAD 1:
	// HTTP rodando

	// AO MESMO TEMPO

	// THREAD 2:
	// baixando segmentos
	go FetchVideo()
	// segment_000 chega
	// ↓
	// player já consegue tocar
	// enquanto isso segment_001 ainda está baixando

	srv := &http.Server{
		Addr:    "localhost:3000",
		Handler: mux,
	}

	fmt.Println("Servidor HTTP iniciado em http://localhost:3000/index.html")

	exec.Command("cmd", "/c", "start", "http://localhost:3000/index.html").Start()

	err := srv.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}

}
