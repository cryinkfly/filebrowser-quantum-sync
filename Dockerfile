# filebrowser-quantum-sync
# Author: Steve Zabka
# Author-URL: https://cryinkfly.com
# License:  Apache-2.0
# 
# Version: 1.0.0

# Build Stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o filebrowser-quantum-sync

# Runtime Stage
FROM scratch
WORKDIR /
COPY --from=builder /app/filebrowser-quantum-sync /filebrowser-quantum-sync
ENTRYPOINT ["/filebrowser-quantum-sync"]
