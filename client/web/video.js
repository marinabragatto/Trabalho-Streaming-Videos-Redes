const video = document.getElementById("video-player");
const bufferVideo = document.getElementById("video-buffer");

const params = new URLSearchParams(window.location.search);
const videoId = params.get("id");

let segments = [];
let current = 0;

// carrega o manifesto
async function loadManifest() {
    const response = await fetch(`/manifest?id=${videoId}`);
    const manifest = await response.json();

    segments = manifest.segments.map(
        segment => `/stream?id=${videoId}&segment=${segment}`
    );

    playSegment(current);
}

//carrega proximo segmento antes do atual acabar
function preloadNextSegment() {
    const next = current + 1;

    if (next < segments.length) {
        bufferVideo.src = segments[next];
        bufferVideo.load();

        console.log("Pré-carregando:", segments[next]);
    }
}

// toca um segmento
function playSegment(index) {
    if (index >= segments.length) {
        console.log("Fim do vídeo");
        return;
    }

    video.src = segments[index];
    video.load();

    video.oncanplay = () => {
        video.play();
        preloadNextSegment();
    };

    console.log("Reproduzindo:", segments[index]);
}

video.addEventListener("ended", () => {
    current++;

    if (current >= segments.length) {
        console.log("Fim do vídeo");
        return;
    }

    video.src = bufferVideo.src;
    video.load();

    video.oncanplay = () => {
        video.play();
        preloadNextSegment();
    };
});

loadManifest();


// const video = document.getElementById("video-player");

// const params = new URLSearchParams(window.location.search);
// const videoId = params.get("id");

// console.log("videoId:", videoId);
// let segments = [];
// let current = 0;

// // carrega o manifesto
// async function loadManifest() {

//     const response = await fetch(`/manifest?id=${videoId}`);
//     console.log("manifest:", videoId);

//     const manifest = await response.json();

//     segments = manifest.segments.map(
//         segment =>  `/stream?id=${videoId}&segment=${segment}`
//     );

//     console.log("Manifest carregado:", segments);

//     playSegment(current);
// }

// // toca um segmento
// function playSegment(index) {

//     if (index >= segments.length) {
//         console.log("Fim do vídeo");
//         return;
//     }

//     video.src = segments[index];

//     video.load();

//     video.play();

//     console.log("Reproduzindo:", segments[index]);
// }

// // quando acabar um segmento → toca o próximo
// video.addEventListener("ended", () => {

//     current++;

//     playSegment(current);
// });

// // inicia lendo o manifesto
// loadManifest();