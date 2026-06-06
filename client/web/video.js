console.log("VIDEO.JS MEDIASOURCE COM AUDIO");

// Pega o player do HTML
const video = document.getElementById("video-player");

// Pega o id da URL. Exemplo: /video?id=2
const params = new URLSearchParams(window.location.search);
const videoId = params.get("id");

// Controle de qualidade
const qualitySelect = document.getElementById("quality-select");

const quality = params.get("quality") || "1080";

qualitySelect.value = quality;
qualitySelect.addEventListener("change", () => {

    params.set("quality", qualitySelect.value);

    window.location.search = params.toString();
});

// MediaSource é a fonte de mídia controlada pelo JavaScript
let mediaSource;

// Buffers separados: um para vídeo e outro para áudio
let videoBuffer;
let audioBuffer;

// Listas de segmentos vindas do manifest.json
let videoSegments = [];
let audioSegments = [];

// Índices dos próximos segmentos a serem carregados
let currentVideoSegment = 0;
let currentAudioSegment = 0;

// Controle para não tentar append enquanto o buffer está ocupado
let appendingVideo = false;
let appendingAudio = false;

// Carrega o manifest.json pelo backend Go
async function loadManifest() {
    const response = await fetch(`/manifest?id=${videoId}`);
    const manifest = await response.json();

    videoSegments = manifest.video;
    audioSegments = manifest.audio;

    console.log("Segmentos de vídeo:", videoSegments);
    console.log("Segmentos de áudio:", audioSegments);

    startMediaSource();
}

// Cria o MediaSource e conecta ele ao <video>
function startMediaSource() {
    mediaSource = new MediaSource();

    // O src do vídeo agora vira blob, não /stream direto
    video.src = URL.createObjectURL(mediaSource);

    mediaSource.addEventListener("sourceopen", () => {
    
        console.log("MediaSource aberto");

        // Codec do vídeo
        const videoCodec = 'video/mp4; codecs="avc1.4d401e"';

        // Codec do áudio AAC
        const audioCodec = 'audio/mp4; codecs="mp4a.40.2"';

        videoBuffer = mediaSource.addSourceBuffer(videoCodec);
        audioBuffer = mediaSource.addSourceBuffer(audioCodec);

        video.addEventListener("timeupdate", () => {
            const buffered = video.buffered.length > 0 ? video.buffered.end(0) : 0;
            const ahead = buffered - video.currentTime;

            // adiciona 12 segundos ao buffer. 
            // como cada chunk tem mais ou menos 4 segundos, adiciona 3 chunks
            // e vai sempre mantendo essa diferenca de 12 segundos a frente, ou seja, sempre
            // que a menos de 12 segundos de video carregado, adiciona outro chunk pra manter
            // a diferenca 
            if (ahead < 12) {
                appendNextVideoSegment();
                appendNextAudioSegment();
            }
        });

        video.addEventListener("waiting", () => {
            appendNextVideoSegment();
            appendNextAudioSegment();
        });

        // Só reseta o flag de "ocupado", não chama o próximo segmento
        videoBuffer.addEventListener("updateend", () => {
            appendingVideo = false;
            tryEndStream();
        });

        audioBuffer.addEventListener("updateend", () => {
            appendingAudio = false;
            tryEndStream();
        });

        videoBuffer.addEventListener("error", (e) => console.error("Erro no videoBuffer:", e));
        audioBuffer.addEventListener("error", (e) => console.error("Erro no audioBuffer:", e));
        video.addEventListener("error", () => console.error("Erro no elemento video:", video.error));

        // Carga inicial — sem isso o vídeo não começa pois o timeupdate
        // só dispara depois que o vídeo já está tocando
        appendNextVideoSegment();
        appendNextAudioSegment();
    });
}

// Adiciona o próximo segmento de vídeo ao buffer de vídeo
async function appendNextVideoSegment() {
    if (appendingVideo) return;
    if (!videoBuffer || videoBuffer.updating) return;

    if (currentVideoSegment >= videoSegments.length) return;

    const segmentName = videoSegments[currentVideoSegment];
    currentVideoSegment++;

    appendingVideo = true;

    const data = await fetchSegment(segmentName);
    videoBuffer.appendBuffer(data);
}

// Adiciona o próximo segmento de áudio ao buffer de áudio
async function appendNextAudioSegment() {
    if (appendingAudio) return;
    if (!audioBuffer || audioBuffer.updating) return;

    if (currentAudioSegment >= audioSegments.length) return;

    const segmentName = audioSegments[currentAudioSegment];
    currentAudioSegment++;

    appendingAudio = true;

    const data = await fetchSegment(segmentName);
    audioBuffer.appendBuffer(data);
}

// Busca um segmento no backend Go.
// O backend usa TCP para buscar no servidor e devolve os bytes ao navegador.
async function fetchSegment(segmentName) {
    const url = 
    `/stream?id=${videoId}` + 
    `&quality=${quality}` +
    `&segment=${segmentName}`;
    console.log("Carregando:", url);

    const response = await fetch(url);

    if (!response.ok) {
        throw new Error("Erro ao carregar segmento: " + segmentName);
    }

    return await response.arrayBuffer();
}

// Só finaliza o MediaSource quando vídeo e áudio acabaram
function tryEndStream() {
    const videoEnded = currentVideoSegment >= videoSegments.length;
    const audioEnded = currentAudioSegment >= audioSegments.length;

    const buffersIdle =
        videoBuffer &&
        audioBuffer &&
        !videoBuffer.updating &&
        !audioBuffer.updating &&
        !appendingVideo &&
        !appendingAudio;

    if (videoEnded && audioEnded && buffersIdle) {
        if (mediaSource.readyState === "open") {
            console.log("Fim do vídeo e do áudio");
            mediaSource.endOfStream();
        }
    }
}

loadManifest();