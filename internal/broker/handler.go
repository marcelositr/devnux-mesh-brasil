package broker

import (
	"log/slog"
	"strings"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"

	meshtastic "github.com/kmpm/meshtastic-protobufs.go/v2/generated"

	"github.com/marcelositr/devnux-mesh-brasil/internal/config"
)

// Handler integra os eventos MQTT com a inspeção opcional de pacotes Meshtastic.
type Handler struct {
	mqtt.HookBase
	logger    *slog.Logger
	processor packetProcessor
}

// NewHandler cria o hook responsável pelo ciclo de vida das mensagens MQTT.
func NewHandler(logger *slog.Logger, cfg config.Config) *Handler {
	return &Handler{
		logger:    logger,
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
	h.logger.Debug("cliente conectado", "client_id", client.ID)
	return nil
}

// OnSubscribe registra a inscrição e preserva os filtros recebidos pelo cliente.
func (h *Handler) OnSubscribe(client *mqtt.Client, packet packets.Packet) packets.Packet {
	h.logger.Debug("inscrição MQTT recebida", "client_id", client.ID, "filters", packet.Filters)
	return packet
}

// OnPublish inspeciona pacotes Meshtastic quando possível, sem filtrar o fluxo MQTT.
// A inspeção e a PSK local servem apenas para observabilidade; o tópico e o payload
// originais são preservados para que o broker encaminhe o tráfego aos assinantes.
func (h *Handler) OnPublish(client *mqtt.Client, packet packets.Packet) (packets.Packet, error) {
	if len(packet.Payload) == 0 {
		h.logger.Debug("publicação MQTT com payload vazio encaminhada", "topic", packet.TopicName, "client_id", client.ID)
		return packet, nil
	}

	// JSON, MapReport, status e outros tópicos não carregam necessariamente
	// um ServiceEnvelope. Eles continuam no fluxo MQTT sem interpretação local.
	if !isProtobufTopic(packet.TopicName) {
		h.logger.Debug("publicação MQTT não protobuf encaminhada", "topic", packet.TopicName, "client_id", client.ID, "payload_bytes", len(packet.Payload))
		return packet, nil
	}

	envelope, err := h.processor.parseServiceEnvelope(packet.Payload)
	if err != nil {
		h.logger.Warn("falha ao inspecionar ServiceEnvelope; publicação preservada", "topic", packet.TopicName, "client_id", client.ID)
		return packet, nil
	}

	if envelope.GetPacket() != nil && envelope.GetPacket().GetDecoded() != nil {
		h.logReceivedMessage(packet.TopicName, client.ID, envelope.GetPacket().GetDecoded())
		return packet, nil
	}

	if err := h.processor.validateServiceEnvelope(envelope); err != nil {
		h.logger.Warn("ServiceEnvelope não pôde ser validado para inspeção; publicação preservada", "topic", packet.TopicName, "client_id", client.ID, "erro", err)
		return packet, nil
	}

	data, err := h.processor.decryptMeshPacket(envelope)
	if err != nil {
		// A chave configurada pode não corresponder ao canal ou ao tipo de
		// criptografia. Isso não pode impedir o encaminhamento do pacote bruto.
		h.logger.Debug("pacote Meshtastic encaminhado sem descriptografia local", "topic", packet.TopicName, "client_id", client.ID, "erro", err)
		return packet, nil
	}
	if data == nil {
		h.logger.Debug("pacote Meshtastic sem dados interpretáveis localmente encaminhado", "topic", packet.TopicName, "client_id", client.ID)
		return packet, nil
	}

	h.logReceivedMessage(packet.TopicName, client.ID, data)
	return packet, nil
}

// isProtobufTopic reconhece tópicos protobuf atuais e o prefixo legado.
// Outros tópicos MQTT continuam sendo encaminhados sem parsing específico.
func isProtobufTopic(topic string) bool {
	parts := strings.Split(strings.Trim(topic, "/"), "/")
	for i := 0; i+1 < len(parts); i++ {
		if parts[i] == "2" && (parts[i+1] == "e" || parts[i+1] == "c") {
			return true
		}
	}
	return false
}

// logReceivedMessage registra metadados, nunca o conteúdo do payload de texto.
func (h *Handler) logReceivedMessage(topic, clientID string, data *meshtastic.Data) {
	if data != nil && int32(data.GetPortnum()) == 1 {
		h.logger.Info("mensagem de texto recebida", "topic", topic, "client_id", clientID, "portnum", int32(data.GetPortnum()), "payload_bytes", len(data.GetPayload()))
		return
	}

	var portnum int32
	if data != nil {
		portnum = int32(data.GetPortnum())
	}

	h.logger.Info("pacote recebido", "topic", topic, "client_id", clientID, "portnum", portnum)
}
