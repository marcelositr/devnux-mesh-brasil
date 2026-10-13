package broker

import (
	"log/slog"
	"strings"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"google.golang.org/protobuf/proto"

	meshtastic "github.com/kmpm/meshtastic-protobufs.go/v2/generated"

	"github.com/marcelositr/devnux-mesh-brasil/internal/config"
)

// Handler integrates MQTT events with optional Meshtastic packet inspection.
type Handler struct {
	mqtt.HookBase
	logger    *slog.Logger
	processor packetProcessor
	zeroHop   bool
}

// NewHandler creates the hook responsible for MQTT message lifecycle events.
func NewHandler(logger *slog.Logger, cfg config.Config) *Handler {
	return &Handler{
		logger:    logger,
		processor: newPacketProcessor(cfg.DefaultPSK),
		zeroHop:   cfg.ZeroHop,
	}
}

func (h *Handler) ID() string {
	return "meshtastic-packet-handler"
}

// Provides limits the hook to MQTT events processed by this broker.
func (h *Handler) Provides(event byte) bool {
	switch event {
	case mqtt.OnConnect, mqtt.OnSubscribe, mqtt.OnPublish:
		return true
	default:
		return false
	}
}

func (h *Handler) Init(any) error {
	return nil
}

// OnConnect records accepted connections without adding authentication.
func (h *Handler) OnConnect(client *mqtt.Client, _ packets.Packet) error {
	h.logger.Debug("cliente conectado", "client_id", client.ID)
	return nil
}

// OnSubscribe records subscriptions and preserves the filters received.
func (h *Handler) OnSubscribe(client *mqtt.Client, packet packets.Packet) packets.Packet {
	h.logger.Debug("inscrição MQTT recebida", "client_id", client.ID, "filters", packet.Filters)
	return packet
}

// OnPublish inspects Meshtastic protobuf packets when possible without filtering
// MQTT traffic. Local inspection failures never prevent the original publication
// from being forwarded. The payload is rewritten only when zero-hop is enabled
// and a valid ServiceEnvelope's hop_limit must be changed.
func (h *Handler) OnPublish(client *mqtt.Client, packet packets.Packet) (packets.Packet, error) {
	if len(packet.Payload) == 0 {
		h.logger.Debug(
			"publicação MQTT com payload vazio encaminhada",
			"topic", packet.TopicName,
			"client_id", client.ID,
		)
		return packet, nil
	}

	// JSON, MapReport, status, and other MQTT topics do not necessarily contain
	// ServiceEnvelope protobufs. Preserve them without local interpretation.
	if !isProtobufTopic(packet.TopicName) {
		h.logger.Debug(
			"publicação MQTT não protobuf encaminhada",
			"topic", packet.TopicName,
			"client_id", client.ID,
			"payload_bytes", len(packet.Payload),
		)
		return packet, nil
	}

	envelope, err := h.processor.parseServiceEnvelope(packet.Payload)
	if err != nil {
		h.logger.Warn(
			"falha ao inspecionar ServiceEnvelope; publicação preservada",
			"topic", packet.TopicName,
			"client_id", client.ID,
			"payload_bytes", len(packet.Payload),
			"erro", err,
		)
		return packet, nil
	}

	if envelope.GetPacket() != nil && envelope.GetPacket().GetDecoded() != nil {
		h.logReceivedMessage(packet.TopicName, client.ID, envelope.GetPacket().GetDecoded())
		return packet, nil
	}

	if err := h.processor.validateServiceEnvelope(envelope); err != nil {
		h.logger.Warn(
			"ServiceEnvelope não pôde ser validado para inspeção; publicação preservada",
			"topic", packet.TopicName,
			"client_id", client.ID,
			"erro", err,
		)
		return packet, nil
	}

	// The public Meshtastic broker limits remaining hops to zero for packets
	// it forwards, preventing MQTT downlink from propagating through the local
	// LoRa mesh. Keep hop_start unchanged to preserve the original information.
	if h.zeroHop && envelope.GetPacket().GetHopLimit() != 0 {
		envelope.GetPacket().HopLimit = 0

		rewritten, marshalErr := proto.Marshal(envelope)
		if marshalErr != nil {
			h.logger.Warn(
				"não foi possível aplicar zero-hop; publicação original preservada",
				"topic", packet.TopicName,
				"client_id", client.ID,
				"erro", marshalErr,
			)
			return packet, nil
		}

		packet.Payload = rewritten
		h.logger.Debug(
			"política zero-hop aplicada",
			"topic", packet.TopicName,
			"client_id", client.ID,
			"hop_start", envelope.GetPacket().GetHopStart(),
		)
	}

	data, err := h.processor.decryptMeshPacket(envelope)
	if err != nil {
		// The local PSK may not match the channel. Inspection failures must not
		// block MQTT forwarding, and are logged without exposing payload data.
		h.logger.Info(
			"pacote Meshtastic recebido; sem descriptografia local",
			"topic", packet.TopicName,
			"client_id", client.ID,
			"payload_bytes", len(packet.Payload),
			"erro", err,
		)
		return packet, nil
	}

	if data == nil {
		h.logger.Info(
			"pacote Meshtastic recebido; sem dados interpretáveis localmente",
			"topic", packet.TopicName,
			"client_id", client.ID,
			"payload_bytes", len(packet.Payload),
		)
		return packet, nil
	}

	h.logReceivedMessage(packet.TopicName, client.ID, data)
	return packet, nil
}

// isProtobufTopic recognizes current Meshtastic protobuf topics and the legacy
// prefix. Other MQTT topics continue to be forwarded without protobuf parsing.
func isProtobufTopic(topic string) bool {
	parts := strings.Split(strings.Trim(topic, "/"), "/")
	for i := 0; i+1 < len(parts); i++ {
		if parts[i] == "2" && (parts[i+1] == "e" || parts[i+1] == "c") {
			return true
		}
	}
	return false
}

// logReceivedMessage records metadata only; it never logs text payload content.
func (h *Handler) logReceivedMessage(topic, clientID string, data *meshtastic.Data) {
	if data != nil && int32(data.GetPortnum()) == 1 {
		h.logger.Info(
			"mensagem de texto recebida",
			"topic", topic,
			"client_id", clientID,
			"portnum", int32(data.GetPortnum()),
			"payload_bytes", len(data.GetPayload()),
		)
		return
	}

	var portnum int32
	if data != nil {
		portnum = int32(data.GetPortnum())
	}
	h.logger.Info("pacote recebido", "topic", topic, "client_id", clientID, "portnum", portnum)
}
