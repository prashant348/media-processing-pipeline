# Media Processing Pipeline

## What is it?
A distributed, cloud-native Video-on-Demand (VoD) orchestration engine built in Go. It manages the lifecycle of a video from raw upload to globally distributed HLS segments.

## What problem does it solves?
It solves the "Clogged Pipe" problem. Standard HTTP downloads (Progressive Download) fail at scale because they consume massive bandwidth and don't adapt to user internet speeds. This project enables Adaptive Bitrate Streaming (ABR) and asynchronous processing to handle thousands of concurrent viewers and massive file uploads without crashing the server.

## How does it solves?
By decoupling the Ingest, Transcoding, and Delivery layers. It uses Go's concurrency (goroutines/channels) to manage FFmpeg workers and leverages object storage (S3) and CDNs to offload the heavy lifting of data delivery.

## **Phase 1: The Ingest Service (Go + S3)**

Don't stream directly from your Go server's disk.

- **The Flow:** User uploads a `.mp4` $\rightarrow$ Go API receives it $\rightarrow$ Go streams the upload directly to "Raw Storage" (e.g., AWS S3 or MinIO).
- **Go's Role:** Use `io.Pipe` and `s3manager.Uploader` to stream the data so you don't load the entire 2GB video into your server's RAM.

---

**Status:** 🚧 Work in Progress - Building Phase 1  
**Last Updated:** 04/07/2026 (dd/mm/yyyy)