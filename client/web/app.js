const video = document.getElementById("video");

let segments = [];
let current = 0;

// carrega o manifesto
async function loadManifest() {

    const response = await fetch("../segments/manifest.json");

    const manifest = await response.json();

    // adiciona /segments/ antes de cada nome
    segments = manifest.segments.map(
        segment => "/segments/" + segment
    );

    console.log("Manifest carregado:", segments);

    playSegment(current);
}

// toca um segmento
function playSegment(index) {

    if (index >= segments.length) {
        console.log("Fim do vídeo");
        return;
    }

    video.src = segments[index];

    video.load();

    video.play();

    console.log("Reproduzindo:", segments[index]);
}

// quando acabar um segmento → toca o próximo
video.addEventListener("ended", () => {

    current++;

    playSegment(current);
});

// inicia lendo o manifesto
loadManifest();