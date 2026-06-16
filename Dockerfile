# ---- build stage ----
FROM golang:1.25-alpine AS build
WORKDIR /src
# Cache module downloads (RAMen has zero third-party deps, but this keeps the
# layer cache friendly if any are added later).
COPY go.mod ./
RUN go mod download
COPY . .
ARG VERSION=docker
RUN CGO_ENABLED=0 go build -trimpath \
    -ldflags "-s -w -X github.com/Rohit-Dnath/RAMen/internal/server.Version=${VERSION}" \
    -o /out/ramen ./cmd/ramen

# ---- runtime stage ----
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/ramen /ramen
# 6379: RESP server (drop-in for REDIS_URL). 8080: web dashboard.
EXPOSE 6379 8080
# Persist snapshots to a mounted volume at /data.
VOLUME ["/data"]
ENV RAMEN_SNAPSHOT_PATH=/data/ramen.snapshot
ENTRYPOINT ["/ramen"]
