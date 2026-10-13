package broker

import (
	"errors"
	"fmt"
	"strings"

	meshtastic "github.com/kmpm/meshtastic-protobufs.go/v2/generated"
	"github.com/marcelositr/devnux-mesh-brasil/internal/crypto"
	"google.golang.org/protobuf/proto"
)

type packetProcessor struct {
	defaultPSK []byte
}

func newPacketProcessor(psk []byte) packetProcessor {
	key := append([]byte(nil), psk...)
	return packetProcessor{defaultPSK: key}
}

func (p packetProcessor) parseServiceEnvelope(payload []byte) (*meshtastic.ServiceEnvelope, error) {
	var envelope meshtastic.ServiceEnvelope
	if err := proto.Unmarshal(payload, &envelope); err != nil {
		return nil, err
	}
	return &envelope, nil
}

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
	if packet.GetId() == 0 {
		return errors.New("packet ID inválido")
	}
	if packet.GetFrom() == 0 {
		return errors.New("origem do pacote inválida")
	}
	if len(packet.GetEncrypted()) == 0 && packet.GetDecoded() == nil {
		return errors.New("payload criptografado ou decodificado ausente")
	}
	return nil
}

func (p packetProcessor) decryptMeshPacket(envelope *meshtastic.ServiceEnvelope) (*meshtastic.Data, error) {
	if envelope == nil || envelope.GetPacket() == nil {
		return nil, errors.New("service envelope ou packet ausente")
	}
	packet := envelope.GetPacket()
	if len(packet.GetEncrypted()) == 0 {
		return nil, errors.New("payload criptografado ausente")
	}

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
