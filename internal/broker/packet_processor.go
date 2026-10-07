package broker

import (
	"errors"
	"fmt"
	"strings"

	meshtastic "github.com/kmpm/meshtastic-protobufs.go/v2/generated"
	"google.golang.org/protobuf/proto"

	"github.com/marcelositr/devnux-mesh-brasil/internal/crypto"
)

// packetProcessor concentra o processamento do protocolo Meshtastic recebido pelo broker.
type packetProcessor struct {
	defaultPSK []byte
}

func newPacketProcessor(defaultPSK []byte) packetProcessor {
	return packetProcessor{defaultPSK: defaultPSK}
}

// parseServiceEnvelope converte o payload MQTT no envelope protobuf usado pelo protocolo.
func (p packetProcessor) parseServiceEnvelope(payload []byte) (*meshtastic.ServiceEnvelope, error) {
	var envelope meshtastic.ServiceEnvelope
	if err := proto.Unmarshal(payload, &envelope); err != nil {
		return nil, err
	}
	return &envelope, nil
}

// validateServiceEnvelope verifica os campos mínimos necessários para processar um pacote.
func (p packetProcessor) validateServiceEnvelope(envelope *meshtastic.ServiceEnvelope) error {
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
func (p packetProcessor) decryptMeshPacket(envelope *meshtastic.ServiceEnvelope) (*meshtastic.Data, error) {
	packet := envelope.GetPacket()
	nonce := crypto.NewNonce(packet.GetFrom(), packet.GetId())

	decrypted, err := crypto.TransformPacket(packet.GetEncrypted(), nonce, p.defaultPSK)
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
