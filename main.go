package main

import (
	"Trabalho-Streaming-Videos-Redes/client"
	"Trabalho-Streaming-Videos-Redes/server"
	"time"
)

func  main()  {
	go server.Serve()    
	time.Sleep(1 * time.Second)
	client.Init()
}