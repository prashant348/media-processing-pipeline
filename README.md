# 🎬 Media Processing Pipeline

A small but practical media processing pipeline built in Go for uploading videos, processing them asynchronously, and exposing HLS-ready output through HTTP APIs.

---

## 🚀 Overview

This project is a backend-first video processing service that takes a raw video upload, stores it in MinIO, creates a background processing job, and uses FFmpeg to generate HLS-compatible segments and playlists.

In its current version, the system is focused on a simple and practical workflow:

1. accept a video upload through an API endpoint
2. store the source file in object storage
3. queue a transcoding job
4. process the video in the background
5. expose status and streaming endpoints for clients

It is designed as a learning-oriented and extensible pipeline rather than a full production-scale media platform.

---

## ⚡ Core Features

- upload video files through a REST API
- store source videos in MinIO
- process jobs asynchronously using a queue and worker pool
- transcode videos into HLS format with FFmpeg
- generate playlist and segment files such as .m3u8 and .ts
- track job status from pending to processing to completed or failed
- expose job and stream endpoints for API clients such as Postman, Requestly, and curl

---

## 🧠 Project Scope

This version is intentionally focused on the backend pipeline:

- ingesting media through HTTP
- decoupling processing from request handling
- producing HLS output files
- providing API-based status and streaming access

A dedicated frontend/dashboard is planned for future versions, so the current implementation is centered around API-driven workflows.

---

## 🏗️ Architecture

The system follows a simple queue-worker architecture:

```text
API Client
    ↓
Go API Server
    ↓
MinIO Object Storage
    ↓
Job Queue
    ↓
Worker Pool
    ↓
FFmpeg Transcoding
    ↓
HLS Output Files
```

For a more detailed explanation, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

---

## 📦 Tech Stack

### Backend

| Category | Technology | Name | Purpose |
|---|---|---|---|
| Language | <img src="https://skillicons.dev/icons?i=go"/>| Go | Core backend logic and concurrency handling |
| Media Processing | <img src="https://upload.wikimedia.org/wikipedia/commons/4/4b/FFmpeg-Logo.svg" width="200px" /> | FFmpeg | Video transcoding into HLS segments and playlists |
| Storage | <img src="./assets/MinIO-Logo-Color.svg" width="100px" height="50px" /> | MinIO | S3-compatible object storage for uploaded videos |
| HTTP Server | `net/http` | net/http | API handling and routing |
| Containerization | <img src="https://skillicons.dev/icons?i=docker"/> | Docker | Local infrastructure setup |
| Scripting | <img src="https://skillicons.dev/icons?i=python"/> | Python | Utility scripts (testing, automation, cleanup) |
| Scripting | <img src="https://skillicons.dev/icons?i=bash"/> | Bash | CLI automation and environment setup |

### Testing

- unit tests
- integration tests with MinIO
- end-to-end pipeline tests

---

## 🔄 Current Workflow

1. upload a video via `POST /api/upload`
2. store the file in MinIO
3. create a processing job
4. process the video in the background
5. check progress through `GET /api/status/{job_id}`
6. access HLS output through `GET /api/stream/{video_id}/{file_name}`

---

## 🛠️ Development Setup

```bash
docker compose up --build -d # spin up local infrastructure
go run cmd/api/main.go # start the Go server
```

---

## 🔐 API Endpoints

- `POST /api/upload`
- `GET /api/status/{job_id}`
- `GET /api/job/{job_id}`
- `GET /api/jobs`
- `GET /api/stream/{video_id}/{file_name}`

For request and response details, see [docs/API.md](docs/API.md).

---

## 🚧 Future Improvements

- better retry and failure handling
- progress reporting for long-running jobs
- authentication and rate limiting
- persistent job storage
- support for more advanced HLS workflows
- frontend/dashboard integration in a later version

---

## ⚠️ Notes

- CORS is enabled for local development
- authentication is not implemented yet
- the current version is focused on backend processing and API usage

---

## 📚 Documentation

- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
- [docs/API.md](docs/API.md)

