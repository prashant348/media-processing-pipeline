import { Component } from "../libs/Component";
import { State } from "../libs/State";
import { navigate } from "../main";
import { startPolling } from "../helpers/jobStatus";

interface UploadResponse {
    video_id: string
};

export const url = "http://localhost:8080";

export default class HomePage extends Component {

    render(): string {
        return `
            <h1>Home Page</h1>
            <input id="videoInput" type="file" />
            <button id="uploadBtn">Upload Video</button>
            <p id="statusBox"></p>
            <button id="watchBtn" style="display: none;">Watch Video</button>
        `
    }

    mount(): void {
        const uploadBtn = document.querySelector<HTMLButtonElement>("#uploadBtn")!;
        const videoInput = document.querySelector<HTMLInputElement>("#videoInput")!;
        const statusBox = document.querySelector<HTMLParagraphElement>("#statusBox")!;  
        const watchBtn = document.querySelector<HTMLButtonElement>("#watchBtn")!;

        const statusState = new State<string>("", (val) => {
            if (uploadBtn) uploadBtn.disabled = val? (val ==="completed"? false: true): false;
            if (statusBox) statusBox.innerText = val;
            if (watchBtn) watchBtn.style.display = val === "completed" ? "block" : "none";
        });

        let videoId: string;

        watchBtn.addEventListener("click", () => navigate(`/video/${videoId}`));

        uploadBtn.addEventListener("click", async (e) => {
            e.preventDefault();

            console.log("upload button clicked!");

            if (videoInput.files?.length === 0) return alert("Please select a file");

            videoId = ""

            const file = videoInput.files?.[0]!;
            const formData = new FormData();
            formData.append("file", file);
            videoInput.value = "";

            try {
                statusState.val = "Uploading...";
                const res = await fetch(`${url}/api/upload`, {
                    method: "POST",
                    body: formData
                });

                if (!res.ok) {
                    statusState.val = "Upload failed!"
                    const errTxt = await res.text();
                    throw new Error(errTxt || "Upload fail");
                };

                statusState.val = "Uploaded";

                const json: UploadResponse = await res.json();
                videoId = json.video_id;
                console.log("server response: ", videoId);

                startPolling(2, videoId, statusState);

            } catch (err) {
                console.error("Upload fail: ", err);
            }
        });
    }
}