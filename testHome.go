package main

import (
	"Trabalho-Streaming-Videos-Redes/client"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
    ip := client.ReadIP()
    const clients = 20000

    var sucessos int64
    var falhas int64
    start := time.Now()

    var wg sync.WaitGroup

    for i := 0; i < clients; i++ {
        wg.Add(1)

        go func(id int) {
            defer wg.Done()

            _, err := client.DoRequestListVideos(ip)

            if err != nil {
                atomic.AddInt64(&falhas, 1)
                fmt.Printf("Cliente %d falhou: %v\n", id, err)
            } else {
                atomic.AddInt64(&sucessos, 1)
            }
        }(i)
    }

    wg.Wait()
    fmt.Printf("Sucessos: %d\n", sucessos)
    fmt.Printf("Falhas: %d\n", falhas)
    fmt.Printf("Tempo total: %v\n", time.Since(start))
}

