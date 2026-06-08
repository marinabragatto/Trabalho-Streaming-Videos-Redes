# Trabalho de Streaming de Vídeos Sob Demanda
Trabalho cliente-servidor com reprodução de vídeos por demanda.
## Descrição
Este projeto consiste em uma aplicação de streaming de vídeos sob demanda desenvolvida utilizando comunicação via sockets TCP e interface web em HTTP. O sistema utiliza vídeos segmentados, permitindo reprodução contínua no navegador sem a necessidade de baixar o arquivo completo antes da execução. 

## Arquitetura
O trabalho vai ser divido em dois terminais diferentes do servidor e um do cliente.
<p align="center">
    <img src="images/1.jpg" width="600">
</p>

* Cliente TCP solicita manifest.json 
* Servidor TCP envia manifest.json
<p align="center">
    <img src="images/2.jpg" width="600">
</p>

* Cliente TCP (Servidor HTML) envia manifest.json para a aplicação HTML (quando solicitada) e ela solicita os segmentos
* Player no navegador
* Reproduz por partes
<p align="center">
    <img src="images/3.jpg" width="600">
</p>

## Tecnologias Utilizadas
#### Linguagens
* Go (Golang)
* JavaScript
* HTML/CSS
#### Bibliotecas e Recursos
* net e http do Go
* API MediaSource do JavaScript

#### MediaSource (funcionamento):
* Lê um arquivo de inicialização (init-stream0.mp4)
* Múltiplos segmentos .m4s
* Controla por conta própria qual pedaço será transmitido.

##### Comando FFMPEG, responsável por preparar os vídeos
```bash
ffmpeg -i 1.mp4 -vf scale=-2:360 -map 0:v:0 -map 0:a:0 -c:v libx264 -c:a aac -seg_duration 4 -use_timeline 0 -use_template 1 -remove_at_exit 0 -init_seg_name 'init-stream$RepresentationID$.mp4' -media_seg_name 'chunk-stream$RepresentationID$-$Number$.m4s' -f dash segments/manifest.mpd
```
Os vídeos são convertidos pelo FFmpeg utilizando o codec __*H.264*__ para vídeo e __*AAC*__ para áudio. Em seguida, o conteúdo é segmentado no formato __*MPEG-DASH*__ em blocos de aproximadamente 4 segundos, permitindo que o navegador solicite e reproduza apenas os segmentos necessários durante o streaming.

## Funcionalidades Implementadas
* Servidor TCP para envio de segmentos de vídeo e áudio
* Servidor HTTP para comunicação com o navegador
* Reprodução de vídeo sob demanda
* Streaming segmentado utilizando manifestos
* Controle de qualidade do vídeo
* Bufferização contínua utilizando MediaSource
* Suporte a múltiplos clientes simultâneos através de goroutines

## Melhorias futuras 
* Implementação de sincronização durante a troca de qualidade do vídeo, permitindo que a reprodução continue a partir do ponto exato em que estava antes da alteração de resolução, proporcionando uma experiência mais fluida ao usuário
* Realização testes de carga e desempenho para dimensionar a capacidade do sistema, avaliando a quantidade de clientes simultâneos suportados e o comportamento da aplicação sob diferentes condições de uso.

## Como executar
Como o github não permite envio de arquivos longos, colocamos os vídeos disponíveis em uma pasta que precisa ser baixada.
>[Link Drive](https://drive.google.com/drive/folders/1pEeOQyyr_Tj7p_rU8H4X5jHdGcassdZm?usp=sharing)
Baixe as pastas de videos e thumbnails e as coloque no diretório do trabalho

### 1. Para rodar o servidor
Em um terminal, execute:
```bash
go run mainServer.go
```
### 2. Para rodar o cliente
Em um outro terminal, execute:
```bash
go run mainClient.go
```
* Em seguida, digite seu IP, no formato (###.###.###..)
* Abra seu navegador em http://localhost:3000

