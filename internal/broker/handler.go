package broker

import (
	"log/slog"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"

	meshtastic "github.com/kmpm/meshtastic-protobufs.go/v2/generated"

	"github.com/marcelositr/devnux-mesh-brasil/internal/config"
)

// Handler integra os eventos MQTT com o processamento de pacotes do broker.
type Handler struct {
	mqtt.HookBase
	logger    *slog.Logger
	processor packetProcessor
}

// NewHandler cria o hook responsável pelo ciclo de vida das mensagens MQTT.
func NewHandler(logger *slog.Logger, cfg config.Config) *Handler {
	return &Handler{
		logger: logger,
		processor: newPacketProcessor(cfg.DefaultPSK),
	}
}

func (h *Handler) ID() string {
	return "meshtastic-packet-handler"
}

// Provides limita o hook aos eventos MQTT que fazem parte do processamento do broker.
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

// OnConnect registra a conexão aceita pelo broker sem aplicar autenticação adicional.
func (h *Handler) OnConnect(client *mqtt.Client, _ packets.Packet) error {
	// A conexão é aceita sem autenticação adicional nesta etapa do broker.
	h.logger.Debug("cliente conectado", "client_id", client.ID)
	return nil
}

// OnSubscribe registra a inscrição e preserva os filtros recebidos pelo cliente.
func (h *Handler) OnSubscribe(client *mqtt.Client, packet packets.Packet) packets.Packet {
	// As inscrições são aceitas sem restrições adicionais nesta etapa do broker.
	h.logger.Debug("inscrição MQTT recebida", "client_id", client.ID, "filters", packet.Filters)
	return packet
}

// OnPublish valida e interpreta o payload Meshtastic antes de permitir a publicação.
// O pacote original é preservado para que o broker MQTT continue seu fluxo normal.
func (h *Handler) OnPublish(client *mqtt.Client, packet packets.Packet) (packets.Packet, error) {
	if len(packet.Payload) == 0 {
		h.logger.Warn("payload vazio recebido", "topic", packet.TopicName, "client_id", client.ID)
		return packet, packets.CodeSuccessIgnore
	}

	envelope, err := h.processor.parseServiceEnvelope(packet.Payload)
	if err != nil {
		h.logger.Warn("falha ao decodificar protobuf", "topic", packet.TopicName, "client_id", client.ID)
		return packet, packets.CodeSuccessIgnore
	}

	if err := h.processor.validateServiceEnvelope(envelope); err != nil {
		h.logger.Warn("service envelope inválido", "topic", packet.TopicName, "client_id", client.ID, "erro", err)
		return packet, packets.CodeSuccessIgnore
	}

	data, err := h.processor.decryptMeshPacket(envelope)
	if err != nil {
		h.logger.Warn("não foi possível interpretar o pacote criptografado", "topic", packet.TopicName, "client_id", client.ID, "erro", err)
		return packet, packets.CodeSuccessIgnore
	}

	h.logReceivedMessage(packet.TopicName, client.ID, data)
	return packet, nil
}

// logReceivedMessage registra mensagens de texto separadamente dos demais portnums.
func (h *Handler) logReceivedMessage(topic, clientID string, data *meshtastic.Data) {
	if data != nil && int32(data.GetPortnum()) == 1 {
		h.logger.Info("mensagem de texto recebida", "topic", topic, "client_id", clientID, "message", string(data.GetPayload()))
		return
	}

	var portnum int32
	if data != nil {
		portnum = int32(data.GetPortnum())
	}

	h.logger.Info("pacote recebido", "topic", topic, "client_id", clientID, "portnum", portnum)
}
