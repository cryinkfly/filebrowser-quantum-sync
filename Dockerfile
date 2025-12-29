# filebrowser-quantum-sync
# Author: Steve Zabka
# Author-URL: https://cryinkfly.com
# License:  Apache-2.0
# 
# Version: 1.0.0

FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -o filebrowser-quantum-sync

FROM scratch
WORKDIR /
COPY --from=builder /app/filebrowser-quantum-sync /filebrowser-quantum-sync
ENTRYPOINT ["/filebrowser-quantum-sync"]
