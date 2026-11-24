# Build stage
FROM golang:1.27-alpine AS builder

RUN apk add --no-cache git ca-certificates

RUN go install github.com/swaggo/swag/cmd/swag@latest    

WORKDIR /build

COPY src/go.mod src/go.sum ./

RUN go mod download

COPY src/ ./

RUN swag init -g main.go --output ./docs --parseDependency

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \ 
    -o /build/bin/api_gateway .

# Run stage
FROM alpine:latest

RUN apk add --no-cache busybox-extras

WORKDIR /api_gateway

COPY --from=builder /build/bin/api_gateway /usr/local/bin/api_gateway
COPY --from=builder /build/docs ./docs

COPY configs/ /etc/api_gateway/configs/

CMD ["api_gateway"]
