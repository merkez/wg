FROM golang:1.15-alpine AS builder
WORKDIR /wg
RUN apk add gcc g++ --no-cache
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o app -a -ldflags '-w -extldflags "-static"'  /wg/grpc/server/main.go

FROM alpine
WORKDIR /app
RUN apk update && apk add sudo && apk add iptables && apk add -U wireguard-tools
COPY --from=builder /wg/config/config.yml /app/config.yml
COPY --from=builder /wg/app /app/app
ENTRYPOINT ["/app/app"]