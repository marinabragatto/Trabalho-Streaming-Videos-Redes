package main

import (
	"Trabalho-Streaming-Videos-Redes/client"
	"Trabalho-Streaming-Videos-Redes/server"
)

func  main()  {
	server.Serve()
	client.Connect_to_server()
}