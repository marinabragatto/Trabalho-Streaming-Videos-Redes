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

//troca a qualidade e salva alguns atributos com informacoes do video anterior 
qualitySelect.addEventListener("change", () => {
    console.log("trocando de qualidade");

    localStorage.setItem("lastVideoId", videoId); 
    localStorage.setItem("lastVideoTime", video.currentTime);
    localStorage.setItem("wasPlaying", !video.paused);

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

// Controle para não tentar update enquanto esta restaurando o video 
let restoring = false;


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

    // O async permite usar await dentro dessa funcao
    mediaSource.addEventListener("sourceopen", async () => {
    
        console.log("MediaSource aberto");

        // Codec do vídeo
        const videoCodec = 'video/mp4; codecs="avc1.4d401e"';

        // Codec do áudio AAC
        const audioCodec = 'audio/mp4; codecs="mp4a.40.2"';

        videoBuffer = mediaSource.addSourceBuffer(videoCodec);
        audioBuffer = mediaSource.addSourceBuffer(audioCodec);        

        // entra nessa funcao sempre que o vídeo avanca no tempo
        video.addEventListener("timeupdate", () => {
            console.log("ADD EVENT LISTENER TIMEUPDATE");

            if(restoring) return; //nao permite carregar seguimentos enquanto esta restaurando

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

        // esse evento acontece quando o vídeo trava esperando dados
        video.addEventListener("waiting", () => {
            console.log("ADD EVENT LISTENER WAITING");

            if(restoring) return; //nao permite carregar seguimentos enquanto esta restaurando

            appendNextVideoSegment();
            appendNextAudioSegment();
        });

        videoBuffer.addEventListener("error", (e) => console.error("Erro no videoBuffer:", e));
        audioBuffer.addEventListener("error", (e) => console.error("Erro no audioBuffer:", e));
        video.addEventListener("error", () => console.error("Erro no elemento video:", video.error));

        

        // busca informacoes do video anterior como id, tempo em que estava e se estava tocando
        const savedVideoId = localStorage.getItem("lastVideoId");
        const savedTime = Number(localStorage.getItem("lastVideoTime"));
        const wasPlaying = localStorage.getItem("wasPlaying") === "true";       

        //se o evento foi apenas uma troca de qualidade e nao uma troca de video
        //entra nessa funcao pra carregar os seguimentos ate o tempo em que estava 
        //antes da troca
        if (savedVideoId === videoId && savedTime > 0) {
            restoring = true;

            console.log("CARREGANDO INIT");

            await appendNextVideoSegment(); // init do vídeo
            await appendNextAudioSegment(); // init do audio

            //calculo estimado para descobrir qual seguimentos estava tocando
            const segmentDuration = 4;
            const segmentIndex = Math.floor(savedTime / segmentDuration);

            currentVideoSegment = segmentIndex;
            currentAudioSegment = segmentIndex; 
            

            //carrega os 4 proximos seguimentos a partir do definido pelo index
            //"carga inicial"
            for (let i = 0; i < 6; i++) {
                console.log("CARREGANDO OUTROS SEGUIMENTOS");

                await appendNextVideoSegment();
                await appendNextAudioSegment();
            }

            
            // Cria uma Promise para esperar o seek terminar.
            // O seek é o ato de mudar video.currentTime para outro ponto.
            const seekPromise = new Promise((resolve) => {
                // Escuta o evento "seeked".
                // Esse evento acontece quando o navegador terminou de ir
                // para o tempo solicitado.
                video.addEventListener("seeked", function onSeeked() {
                    console.log("Seeked para:", video.currentTime);
                    
                    video.removeEventListener("seeked", onSeeked);

                    // Finaliza a Promise e libera o await seekPromise abaixo.
                    resolve();
                });
            });

            video.currentTime = savedTime;

            //espera o evento seeked acontecer
            await seekPromise;

            if (wasPlaying) {
                await video.play();
            }

            localStorage.removeItem("lastVideoTime");
            localStorage.removeItem("wasPlaying");

            restoring = false;
        }        
        else {
            // Caso seja um video novo:
            // Carga inicial — sem isso o vídeo não começa pois o timeupdate
            // só dispara depois que o vídeo já está tocando
            await appendNextVideoSegment();
            await appendNextAudioSegment();
        }
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

    // Cria uma Promise para esperar o appendBuffer terminar.
    return new Promise((resolve) => {
        videoBuffer.addEventListener("updateend", function handler() {
            videoBuffer.removeEventListener("updateend", handler);

            appendingVideo = false;
            tryEndStream();

            resolve();
        });

        videoBuffer.appendBuffer(data);
    });
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

    // Cria uma Promise para esperar o appendBuffer terminar.
    return new Promise((resolve) => {
        audioBuffer.addEventListener("updateend", function handler() {
            audioBuffer.removeEventListener("updateend", handler);

            appendingAudio = false;
            tryEndStream();

            resolve();
        });

        audioBuffer.appendBuffer(data);
    });
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