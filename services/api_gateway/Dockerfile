FROM golang:1.27.4-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/gateway ./services/api_gateway/cmd/gateway

FROM build AS test
CMD ["go", "test", "-count=1", "-v", "./..."]

FROM debian:bookworm-slim AS runtime
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=build /out/gateway /app/gateway
USER 10001:10001
EXPOSE 8080
ENTRYPOINT ["/app/gateway"]
CMD ["serve"]
