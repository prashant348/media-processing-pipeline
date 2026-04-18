import { url } from "../pages/Home";
import { State } from "../libs/State";

interface StatusResponse {
    job_id: string
    status: string
}

export async function fetchJobStatus(
    jobId: string
): Promise<string> {
    const res = await fetch(`${url}/api/status/${jobId}`);

    if (!res.ok) {
        const errTxt = await res.text();
        throw new Error(errTxt || "Failed to fetch job status");
    };

    const json: StatusResponse = await res.json();
    console.log("status: ", json.status)
    return json.status
}

export async function startPolling(
    seconds: number, 
    videoId: string,
    statusState?: State<string>
) {
    const interval = setInterval(async () => {
        const status = await fetchJobStatus(videoId);
        if (statusState) statusState.val = status;
        if (status === "completed" || status === "failed") {
            clearInterval(interval);
        }
    }, seconds * 1000);
}