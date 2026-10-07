package broker

import (
	"bytes"
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

type Handler struct {
	mqtt.HookBase
	logger *slog.Logger
	config config.Config
}

func NewHandler(logger *slog.Logger, cfg config.Config) *Handler {
	return &Handler{logger: logger, config: cfg}
}

func (h *Handler) ID() string {
	return "meshtastic-packet-handler"
}

func (h *Handler) Provides(event byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnSubscribe,
		mqtt.OnPublish,
	}, []byte{event})
}

func (h *Handler) Init(any) error {
	return nil
}

func (h *Handler) OnConnect(client *mqtt.Client, _ packets.Packet) error {
	// O projeto de referência permite a conexão sem autenticação adicional.
	h.logger.Debug("cliente conectado", "client_id", client.ID)
	return nil
}

func (h *Handler) OnSubscribe(client *mqtt.Client, packet packets.Packet) packets.Packet {
	// O projeto de referência permite todas as inscrições.
	h.logger.Debug("inscrição MQTT recebida", "client_id", client.ID, "filters", packet.Filters)
	return packet
}

func (h *Handler) OnPublish(client *mqtt.Client, packet packets.Packet) (packets.Packet, error) {
	if len(packet.Payload) == 0 {
		h.logger.Warn("payload vazio recebido", "topic", packet.TopicName, "client_id", client.ID)
		return packet, packets.CodeSuccessIgnore
	}

	var envelope meshtastic.ServiceEnvelope
	if err := proto.Unmarshal(packet.Payload, &envelope); err != nil {
		h.logger.Warn("falha ao decodificar protobuf", "topic", packet.TopicName, "client_id", client.ID)
		return packet, packets.CodeSuccessIgnore
	}

	if err := validateServiceEnvelope(&envelope); err != nil {
		h.logger.Warn("service envelope inválido", "topic", packet.TopicName, "client_id", client.ID, "erro", err)
		return packet, packets.CodeSuccessIgnore
	}

	data, err := h.decryptMeshPacket(&envelope)
	if err != nil {
		// No projeto de referência o bloqueio por falha de descriptografia está comentado.
		h.logger.Debug("não foi possível interpretar o pacote criptografado", "erro", err)
	}

	h.logReceivedMessage(packet.TopicName, client.ID, data)
	return packet, nil
}

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
		return nil, errors.New("payload Meshtastic inválido")
	}

	return &data, nil
}

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
