package broker

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"google.golang.org/protobuf/proto"

	meshtastic "github.com/kmpm/meshtastic-protobufs.go/v2/generated"

	"github.com/marcelositr/devnux-mesh-brasil/internal/config"
	"github.com/marcelositr/devnux-mesh-brasil/internal/crypto"
)

// Handler integra o ciclo de vida MQTT com o processamento dos pacotes Meshtastic.
type Handler struct {
	mqtt.HookBase
	logger *slog.Logger
	config config.Config
}

// NewHandler cria o hook responsável pelo processamento de mensagens MQTT.
func NewHandler(logger *slog.Logger, cfg config.Config) *Handler {
	return &Handler{logger: logger, config: cfg}
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
	// O projeto de referência permite a conexão sem autenticação adicional.
	h.logger.Debug("cliente conectado", "client_id", client.ID)
	return nil
}

// OnSubscribe registra a inscrição e preserva os filtros recebidos pelo cliente.
func (h *Handler) OnSubscribe(client *mqtt.Client, packet packets.Packet) packets.Packet {
	// O projeto de referência permite todas as inscrições.
	h.logger.Debug("inscrição MQTT recebida", "client_id", client.ID, "filters", packet.Filters)
	return packet
}

// OnPublish valida e interpreta o envelope Meshtastic antes de permitir a publicação.
// O pacote original é preservado para que o broker MQTT continue seu fluxo normal.
func (h *Handler) OnPublish(client *mqtt.Client, packet packets.Packet) (packets.Packet, error) {
	if len(packet.Payload) == 0 {
		h.logger.Warn("payload vazio recebido", "topic", packet.TopicName, "client_id", client.ID)
		return packet, packets.CodeSuccessIgnore
	}

	envelope, err := parseServiceEnvelope(packet.Payload)
	if err != nil {
		h.logger.Warn("falha ao decodificar protobuf", "topic", packet.TopicName, "client_id", client.ID)
		return packet, packets.CodeSuccessIgnore
	}

	if err := validateServiceEnvelope(&envelope); err != nil {
		h.logger.Warn("service envelope inválido", "topic", packet.TopicName, "client_id", client.ID, "erro", err)
		return packet, packets.CodeSuccessIgnore
	}

	data, err := h.decryptMeshPacket(&envelope)
	if err != nil {
		h.logger.Warn("não foi possível interpretar o pacote criptografado", "topic", packet.TopicName, "client_id", client.ID, "erro", err)
		return packet, packets.CodeSuccessIgnore
	}

	h.logReceivedMessage(packet.TopicName, client.ID, data)
	return packet, nil
}

// parseServiceEnvelope converte o payload MQTT no envelope protobuf usado pelo protocolo.
func parseServiceEnvelope(payload []byte) (*meshtastic.ServiceEnvelope, error) {
	var envelope meshtastic.ServiceEnvelope
	if err := proto.Unmarshal(payload, &envelope); err != nil {
		return nil, err
	}
	return &envelope, nil
}

// validateServiceEnvelope verifica os campos mínimos necessários para processar um pacote.
func validateServiceEnvelope(envelope *meshtastic.ServiceEnvelope) error {
	if envelope == nil {
		return errors.New("service envelope nulo")
	}
	if strings.TrimSpace(envelope.GetChannelId()) == "" {
		return errors.New("channel ID vazio")
	}
	if strings.TrimSpace(envelope.GetGatewayId()) == "" {
		return errors.New("gateway ID vazio")
	}

	packet := envelope.GetPacket()
	if packet == nil {
		return errors.New("packet ausente")
	}
	if packet.GetId() < 1 {
		return errors.New("packet ID inválido")
	}
	if packet.GetFrom() < 1 {
		return errors.New("origem do pacote inválida")
	}
	if len(packet.GetEncrypted()) < 1 {
		return errors.New("payload criptografado ausente")
	}
	if packet.GetDecoded() != nil {
		return errors.New("packet já possui payload decodificado")
	}

	return nil
}

// decryptMeshPacket descriptografa o pacote e interpreta o conteúdo como Data.
func (h *Handler) decryptMeshPacket(envelope *meshtastic.ServiceEnvelope) (*meshtastic.Data, error) {
	packet := envelope.GetPacket()
	nonce := crypto.NewNonce(packet.GetFrom(), packet.GetId())

	decrypted, err := crypto.TransformPacket(packet.GetEncrypted(), nonce, h.config.DefaultPSK)
	if err != nil {
		return nil, err
	}

	var data meshtastic.Data
	if err := proto.Unmarshal(decrypted, &data); err != nil {
		return nil, fmt.Errorf("decodificar data: %w", err)
	}

	if int32(data.GetPortnum()) <= 0 || len(data.GetPayload()) == 0 {
		return nil, nil
	}

	return &data, nil
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
