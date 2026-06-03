console.log("INDEX.JS CARREGOU");

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
    // console.log("status:", resposta.status);

    // const texto = await resposta.text();

    // console.log("resposta:", texto);

    // return JSON.parse(texto);

    // const dados = await resposta.json();
    // return dados

// //     // Por enquanto, vamos simular o que o Go enviaria:
//     return [
//         { id: "1", nome: "Toy Story 5", capa: "thumbnails/trailer1.jpg" },
//         { id: "2", nome: "Vingadores", capa: "thumbnails/trailer2.jpg" },
//         { id: "3", nome: "Matrix", capa: "thumbnails/trailer3.jpg" }
//     ];
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
                <a href="/video" class="thumbnail">
                    <img class="thumbnail-image" src="${video.thumbnail}"
                    width="250"
                    height="150"
                    >
                </a>
                    <div class="video-description-section"> 
                    <div class="video-details">
                        <a href="/video"  class="video-title">  ${video.nome}</a>
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



