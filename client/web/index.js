console.log("INDEX.JS CARREGOU");



// console.log("videoId:", videoId);

async function buscarCatalogo() {
    console.log("buscar catalogo");

    // Quando o seu Go estiver pronto, você usaria isso:
    
    try {
        const resposta = await fetch('/api/catalogo');

        console.log("status:", resposta.status);

        const texto = await resposta.text();

        console.log("texto:", texto);

        return JSON.parse(texto);

    } catch (err) {
        console.error("ERRO NO FETCH:", err);
        throw err;
    }
  
}

async function montarTelaInicial() {
    const listaDeVideos = await buscarCatalogo();

    const secaoEmAlta = document.getElementById("lista-em-alta")
    const secaoFavoritos = document.getElementById('lista-favoritos')

    // Limpa o conteúdo por segurança
    secaoEmAlta.innerHTML = '';
    secaoFavoritos.innerHTML = '';

    listaDeVideos.forEach(video => {
      
        const htmlDoVideo = `
            <article class="video-container">
                <a href="/video?id=${video.id}" class="thumbnail">
                    <img class="thumbnail-image" src="${video.thumbnail}"
                    width="250"
                    height="150"
                    >
                </a>
                    <div class="video-description-section"> 
                    <div class="video-details">
                        <a href="/video?id=${video.id}""  class="video-title">  ${video.nome}</a>
                    </div>
                    </div>
        </article>
        `;


        secaoEmAlta.innerHTML += htmlDoVideo;     

        secaoFavoritos.innerHTML += htmlDoVideo;
    });
}

console.log("antes");
montarTelaInicial();
console.log("depois");

// buscarCatalogo();



