package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"Trabalho-Streaming-Videos-Redes/client"
)

type Manifest struct {
	Video []string `json:"video"`
	Audio []string `json:"audio"`
}

var sucessos int64
var falhas int64

func SimularCliente(id int, ip string) {
	
	fmt.Printf("Cliente %d iniciou\n", id)

	manifestBytes, err := client.DoRequestGetManifest(
		1,
		ip,
	)

	if err != nil {
		fmt.Printf("Cliente %d erro manifesto: %v\n", id, err)
		atomic.AddInt64(&falhas, 1)
		return
	}

	var manifest Manifest

	err = json.Unmarshal(manifestBytes, &manifest)
	if err != nil {
		fmt.Printf("Cliente %d erro decode manifesto: %v\n", id, err)
		atomic.AddInt64(&falhas, 1)
		return
	}

	fmt.Printf(
		"Cliente %d recebeu %d segmentos de vídeo e %d de áudio\n",
		id,
		len(manifest.Video),
		len(manifest.Audio),
	)

	// init + segmentos
	for i := 0; i < len(manifest.Video); i++ {
		fmt.Printf("Baixando VIDEO %d: %s\n", i, manifest.Video[i])
		_, err := client.DoRequestGetSegment(
			1,
			1080,
			manifest.Video[i],
			ip,
		)

		if err != nil {
			fmt.Printf(
				"Cliente %d erro segmento vídeo %s: %v\n",
				id,
				manifest.Video[i],
				err,
			)
			atomic.AddInt64(&falhas, 1)
			return
		}

		if i < len(manifest.Audio) {
			fmt.Printf("Baixando AUDIO %d: %s\n", i, manifest.Audio[i])
			_, err := client.DoRequestGetSegment(
				1,
				1080,
				manifest.Audio[i],
				ip,
			)

			if err != nil {
				fmt.Printf(
					"Cliente %d erro segmento áudio %s: %v\n",
					id,
					manifest.Audio[i],
					err,
				)
				atomic.AddInt64(&falhas, 1)
				return
			}
			fmt.Printf("AUDIO %d OK\n", i)
		}
	}

	fmt.Printf("Cliente %d terminou\n", id)

	atomic.AddInt64(&sucessos, 1)
}

func main() {
	ip := client.ReadIP()
	const usuarios = 200

	inicio := time.Now()

	var wg sync.WaitGroup

	for i := 0; i < usuarios; i++ {

		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			SimularCliente(id,ip)
		}(i)
	}

	wg.Wait()

	fmt.Println()
	fmt.Println("===== RESULTADO =====")
	fmt.Printf("Usuários: %d\n", usuarios)
	fmt.Printf("Sucessos: %d\n", sucessos)
	fmt.Printf("Falhas: %d\n", falhas)
	fmt.Printf("Tempo total: %v\n", time.Since(inicio))
}
