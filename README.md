# Trabalho de Streaming de Vídeos Sob Demanda
O trabalho vai ser divido em dois terminais diferentes (possivelmente, para separar o cliente do servidor), mas para facilitar o desenvolvimento, coloquei ambos na main. 
### 3. Baixe os vídeos das pastas do drive e coloque no root do repositório:
Como o github não permite envio de arquivos longos, colocamos os vídeos disponíveis em uma pasta que precisa ser baixada.
>[Link Drive](https://drive.google.com/drive/folders/1pEeOQyyr_Tj7p_rU8H4X5jHdGcassdZm?usp=drive_link)
### 1. Rode a main.go:
Em um terminal, execute a main.
```bash
go run main.go
```

### Explicação das alterações
Usei um programa chamado ffmpeg para seguimentar os vídeos, em pedaços menores de aproximadamente 5 segundos cada (isso depende da janela de cada frame, algo relacionado aos bytes de vídeos que não convem mencionar)

```bash
ffmpeg -i trailer2.mp4 -c:v libx264 -c:a aac -f segment -segment_time 5 -reset_timestamps 1 segment_%03d.mp4
 ```

Além desses segmentos, é necessário também um cabeçalho, *manifest.json* que indica para o cliente a quantidade de segmentos que deve ser exibida.

### Funcionamento
Basicamente o cliente chama 
```bash
go fecthVideo()
```
Esta função, por sua vez, rquisita do servidor primeiramente o manifesto, para que possa ler os segmentos, e a apartir disso requisita os segmentos, um a um para que possa rodá-los um seguido do outro.

O *go*, como já explicado no código, permite que (segment_000 chegue -> player já consiga tocar e enquanto isso -> segment_001 ainda está baixando)

### Problemas 
* Os cortes entre os vídeos ficam pouco sutis

### Arquitetura 
Cliente TCP solicita manifest.json 
↓
Servidor TCP envia manifest.json
↓
Cliente TCP salva manifest.json em em client/segments e usa para saber nome dos seguimentos e solicitá-los
↓
HTTP local
↓
player no navegador
↓
reproduz por partes