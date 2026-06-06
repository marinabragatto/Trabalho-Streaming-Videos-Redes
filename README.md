# Trabalho de Streaming de Vídeos Sob Demanda
O trabalho vai ser divido em dois terminais diferentes (possivelmente, para separar o cliente do servidor), mas para facilitar o desenvolvimento, coloquei ambos na main. 
# 3. Baixe os vídeos das pastas do drive  !NOVAMENTE! e coloque no root do repositório:
## A pasta foi atualizada para funcionar com o novo código!

Como o github não permite envio de arquivos longos, colocamos os vídeos disponíveis em uma pasta que precisa ser baixada.
>[Link Drive](https://drive.google.com/drive/folders/1pEeOQyyr_Tj7p_rU8H4X5jHdGcassdZm?usp=sharing)
### 1. Rode a main.go:
Em um terminal, execute a main.
```bash
go run main.go
```

### Problemas:
* Mudar a qualidade do vídeo e continuar de onde parou, ainda não é possível
* Avançar o vídeo e aparecer na tela o tempo total do vídeo (teria que criar esse atributo no metadata.json do servidor) ainda não é 100% possível
* Interface mais bonita
### Funcionamento
Basicamente o cliente chama 
```bash
go fecthVideo()
```
Esta função, por sua vez, requisita do servidor primeiramente o manifesto, para que possa ler os segmentos, e a apartir disso requisita os segmentos, um a um para que possa rodá-los um seguido do outro.


### Agora utilizamos a API MediaSource para contolar o fluxo dos vídeos, ela lê:
* um arquivo de inicialização (init-stream0.mp4)
* múltiplos segmentos .m4s
E controla por conta própria qual pedaço será transmitido.

### Mesmo comando ffmpeg, só que agora com a flag para definir a qualidade *-vf scale=-2:360*
```bash
ffmpeg -i 1.mp4 -vf scale=-2:360 -map 0:v:0 -map 0:a:0 -c:v libx264 -c:a aac -seg_duration 4 -use_timeline 0 -use_template 1 -remove_at_exit 0 -init_seg_name 'init-stream$RepresentationID$.mp4' -media_seg_name 'chunk-stream$RepresentationID$-$Number$.m4s' -f dash segments/manifest.mpd
```
*É importante que na pasta que você colocar o trailer pra seguimentá-lo, também crie uma pasta segments, na qual o conteúdo será baixado**


**Troquei os endereços para permitir multiplos acessos na mesma rede**
*Obs:* Para funcionar na sua máquina você precisa trocar o campo que tem o meu IP, no tcp.go, pelo seu IP
Para permitir que outros dispositivos da rede local acessem a plataforma de streaming é necessário liberar as portas utilizadas pelo servidor HTTP e pelo servidor TCP no Firewall do Windows.

Execute os comandos abaixo no PowerShell como administrador:
```bash
New-NetFirewallRule -DisplayName "Go HTTP 3000" -Direction Inbound -Protocol TCP -LocalPort 3000 -Action Allow
```
```bash
New-NetFirewallRule -DisplayName "Go TCP 8080" -Direction Inbound -Protocol TCP -LocalPort 8080 -Action Allow
```


### Arquitetura 
* Cliente TCP solicita manifest.json 
* Servidor TCP envia manifest.json
* Cliente TCP envia manifest.json para a aplicação html (quando solicitada) e ela solicita os segmentos
* HTTP local
* player no navegador
* reproduz por partes
