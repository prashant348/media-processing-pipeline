FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ffmpeg ca-certificates

WORKDIR /app

COPY --from=builder /out/server /app/server

RUN mkdir -p /app/output /app/tmp

EXPOSE 8080

CMD ["/app/server"]