import { Component } from "../libs/Component"
import Hls from "hls.js"
import { fetchJobStatus } from "../helpers/jobStatus"

export default class VideoPage extends Component {
    render(): string {
        return `
            <h1>Video Page</h1>
            <div id="video-container">
                <h1 id="video-status">Checking video status...</h1>
                <div id="player-wrapper" style="display: none;">
                    <video id="hlsPlayer" controls style="width: 100%; max-width: 800px;"></video>
                </div>
            </div>
        `
    }

    async mount(): Promise<void> {
        const statusHeader = document.querySelector<HTMLHeadingElement>("#video-status")!;
        const playerWrapper = document.querySelector<HTMLDivElement>("#player-wrapper")!;
        const videoId = this.params.id;

        const checkAndStart = async () => {
            try {
                const status = await fetchJobStatus(videoId);

                if (status === "completed") {
                    statusHeader.style.display = "none"; // Status hide karo
                    playerWrapper.style.display = "block"; // Player show karo
                    this.initHlsPlayer(videoId); // HLS start karo
                    return true; // Polling stop karne ke liye
                }

                if (status === "failed") {
                    statusHeader.innerText = "❌ Video processing failed!";
                    return true;
                }

                statusHeader.innerText = `⏳ Video is ${status}... Please wait.`;
                return false; // Still processing, keep polling
            } catch (err) {
                
                statusHeader.innerText = "⚠️ Error fetching video status.";
                return true;
            }
        };

        // 1. Pehla check turant karein
        const isDone = await checkAndStart();

        // 2. Agar done nahi hai, toh polling start karein
        if (!isDone) {
            const interval = setInterval(async () => {
                const finished = await checkAndStart();
                if (finished) clearInterval(interval);
            }, 3000); // Har 3 second mein check karein
        }
    }

    private initHlsPlayer(id: string) {
    const video = document.getElementById('hlsPlayer') as HTMLVideoElement;
    const videoSrc = `http://localhost:8080/api/stream/${id}/index.m3u8`;

    if (Hls.isSupported()) {
        const hls = new Hls();
        hls.loadSource(videoSrc);
        hls.attachMedia(video);
    }

}

}