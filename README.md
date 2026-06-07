# Trabalho de Streaming de Vídeos Sob Demanda
Trabalho cliente-servidor com reprodução de vídeos por demanda.
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

Como o github não permite envio de arquivos longos, colocamos os vídeos disponíveis em uma pasta que precisa ser baixada.
>[Link Drive](https://drive.google.com/drive/folders/1pEeOQyyr_Tj7p_rU8H4X5jHdGcassdZm?usp=sharing)

### 1. Rode a main.go:
Em um terminal, execute a main.
```bash
go run main.go
```
### Agora utilizamos a API MediaSource para contolar o fluxo dos vídeos, ela lê:
* Um arquivo de inicialização (init-stream0.mp4)
* Múltiplos segmentos .m4s
E controla por conta própria qual pedaço será transmitido.

### Comando FFMPEG, responsável por preparar os vídeos
```bash
ffmpeg -i 1.mp4 -vf scale=-2:360 -map 0:v:0 -map 0:a:0 -c:v libx264 -c:a aac -seg_duration 4 -use_timeline 0 -use_template 1 -remove_at_exit 0 -init_seg_name 'init-stream$RepresentationID$.mp4' -media_seg_name 'chunk-stream$RepresentationID$-$Number$.m4s' -f dash segments/manifest.mpd
```
Os vídeos são convertidos pelo FFmpeg utilizando o codec __*H.264*__ para vídeo e __*AAC*__ para áudio. Em seguida, o conteúdo é segmentado no formato __*MPEG-DASH*__ em blocos de aproximadamente 4 segundos, permitindo que o navegador solicite e reproduza apenas os segmentos necessários durante o streaming.