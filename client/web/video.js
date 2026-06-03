const video = document.getElementById("video-player");

const params = new URLSearchParams(window.location.search);
const videoId = params.get("id");

console.log("videoId:", videoId);
let segments = [];
let current = 0;

// carrega o manifesto
async function loadManifest() {

    const response = await fetch(`/manifest?id=${videoId}`);
    console.log("manifest:", videoId);

    const manifest = await response.json();

    segments = manifest.segments.map(
        segment =>  `/stream?id=${videoId}&segment=${segment}`
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