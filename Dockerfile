FROM golang:1.27.2-alpine@sha256:f92b6ef800e499660581efdabdf25d9d817a9d124eaf900924f0504e7e27e12d AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/devnux-mesh-broker ./cmd/broker

FROM alpine:3.22@sha256:5291449c3df73caf6ed85e649dec1b9e818b39a5d8c871e97afc13e9cd5e8fa8

RUN addgroup -S devnux && adduser -S -G devnux devnux

COPY --from=build /out/devnux-mesh-broker /usr/local/bin/devnux-mesh-broker

LABEL org.opencontainers.image.title="DevNux Mesh Brasil" \
      org.opencontainers.image.description="Broker MQTT comunitário para integração de redes Meshtastic no Brasil" \
      org.opencontainers.image.source="https://github.com/marcelositr/devnux-mesh-brasil" \
      org.opencontainers.image.licenses="AGPL-3.0"

ENV MQTT_LISTEN_ADDRESS=":8883" \
    MQTT_CERTIFICATE_FILE="/etc/devnux/certs/certificate.pem" \
    MQTT_PRIVATE_KEY_FILE="/etc/devnux/certs/private.key"

EXPOSE 8883

USER devnux

ENTRYPOINT ["/usr/local/bin/devnux-mesh-broker"]
