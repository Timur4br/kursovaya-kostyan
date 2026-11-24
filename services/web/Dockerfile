FROM golang:1.27-alpine AS builder
WORKDIR /build
COPY src/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /web .

FROM alpine:3.23
RUN apk add --no-cache ca-certificates
COPY --from=builder /web /usr/local/bin/web
USER 65534:65534
EXPOSE 8080
CMD ["web"]
