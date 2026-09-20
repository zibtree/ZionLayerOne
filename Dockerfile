# ZionLayer node image — September 2026 development baseline

FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /bin/ziond ./cmd/ziond

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
RUN addgroup -S zion && adduser -S zion -G zion
COPY --from=builder /bin/ziond /usr/local/bin/ziond
RUN mkdir -p /data && chown -R zion:zion /data

USER zion
WORKDIR /data

EXPOSE 8545 9000
ENTRYPOINT ["ziond", "start", "--data-dir", "/data"]
