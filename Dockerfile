FROM golang:1.24.4-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/devnux-mesh-broker ./cmd/broker

FROM alpine:3.22

RUN addgroup -S devnux && adduser -S -G devnux devnux

COPY --from=build /out/devnux-mesh-broker /usr/local/bin/devnux-mesh-broker

ENV MQTT_LISTEN_ADDRESS=":8883" \
    MQTT_CERTIFICATE_FILE="/etc/devnux/certs/certificate.pem" \
    MQTT_PRIVATE_KEY_FILE="/etc/devnux/certs/private.key"

EXPOSE 8883

USER devnux

ENTRYPOINT ["/usr/local/bin/devnux-mesh-broker"]
