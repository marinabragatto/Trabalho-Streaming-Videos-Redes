# Trabalho de Streaming de Vídeos Sob Demanda
O trabalho vai ser divido em dois terminais diferentes (possivelmente, para separar o cliente do servidor), mas para facilitar o desenvolvimento, coloquei ambos na main. 
# 3. Baixe os vídeos das pastas do drive  !NOVAMENTE! e coloque no root do repositório:
## A pasta foi atualizada para funcionar com o novo código!

Como o github não permite envio de arquivos longos, colocamos os vídeos disponíveis em uma pasta que precisa ser baixada.
>[Link Drive](https://drive.google.com/drive/folders/10RTWLbGFuI5Jrn6iFp1iTm4LoauU7FSb?usp=drive_link)
### 1. Rode a main.go:
Em um terminal, execute a main.
```bash
go run main.go
```

### Explicação das alterações
Agora todo arquivo passa apenas pela RAM

### Funcionamento
Basicamente o cliente chama 
```bash
go fecthVideo()
```
Esta função, por sua vez, rquisita do servidor primeiramente o manifesto, para que possa ler os segmentos, e a apartir disso requisita os segmentos, um a um para que possa rodá-los um seguido do outro.

O *go*, como já explicado no código, permite que 
* segment_000 chegue -> player já consiga tocar 
* e enquanto isso -> segment_001 ainda está baixando

### Problemas 
* Os cortes entre os vídeos ficam pouco sutis
* Precisa enviar as thumbnails também na página home

### Arquitetura 
* Cliente TCP solicita manifest.json 
* Servidor TCP envia manifest.json
* Cliente TCP envia manifest.json para a aplicação html (quando solicitada) e ela solicita os segmentos
* HTTP local
* player no navegador
* reproduz por partes