FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/ziond ./cmd/ziond

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/ziond /usr/local/bin/ziond
EXPOSE 8545 9000
ENTRYPOINT ["ziond", "start"]
