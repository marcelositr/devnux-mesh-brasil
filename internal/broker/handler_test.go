package broker

import (
	"bytes"
	"io"
	"log/slog"
	"testing"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	meshtastic "github.com/kmpm/meshtastic-protobufs.go/v2/generated"
	"google.golang.org/protobuf/proto"

	"github.com/marcelositr/devnux-mesh-brasil/internal/config"
	"github.com/marcelositr/devnux-mesh-brasil/internal/crypto"
)

func testHandler() *Handler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewHandler(logger, config.Config{DefaultPSK: crypto.DefaultPSK})
}

func testEnvelope(t *testing.T, data *meshtastic.Data) *meshtastic.ServiceEnvelope {
	t.Helper()

	plain, err := proto.Marshal(data)
	if err != nil {
		t.Fatalf("proto.Marshal() retornou erro: %v", err)
	}

	const (
		from     = uint32(123456789)
		packetID = uint32(987654321)
	)

	encrypted, err := crypto.TransformPacket(plain, crypto.NewNonce(from, packetID), crypto.DefaultPSK)
	if err != nil {
		t.Fatalf("TransformPacket() retornou erro: %v", err)
	}

	return &meshtastic.ServiceEnvelope{
		ChannelId: "test-channel",
		GatewayId: "test-gateway",
		Packet: &meshtastic.MeshPacket{
			Id:        packetID,
			From:      from,
			Encrypted: encrypted,
		},
	}
}

func TestValidateServiceEnvelope(t *testing.T) {
	valid := &meshtastic.ServiceEnvelope{
		ChannelId: "channel",
		GatewayId: "gateway",
		Packet: &meshtastic.MeshPacket{
			Id:        1,
			From:      2,
			Encrypted: []byte{1},
		},
	}

	tests := []struct {
		name    string
		envelope *meshtastic.ServiceEnvelope
		wantErr bool
	}{
		{name: "valid", envelope: valid},
		{name: "nil", envelope: nil, wantErr: true},
		{name: "missing channel", envelope: &meshtastic.ServiceEnvelope{
			GatewayId: "gateway",
			Packet:    valid.Packet,
		}, wantErr: true},
		{name: "missing gateway", envelope: &meshtastic.ServiceEnvelope{
			ChannelId: "channel",
			Packet:    valid.Packet,
		}, wantErr: true},
		{name: "missing packet", envelope: &meshtastic.ServiceEnvelope{
			ChannelId: "channel",
			GatewayId: "gateway",
		}, wantErr: true},
		{name: "invalid packet id", envelope: &meshtastic.ServiceEnvelope{
			ChannelId: "channel",
			GatewayId: "gateway",
			Packet: &meshtastic.MeshPacket{
				From:      2,
				Encrypted: []byte{1},
			},
		}, wantErr: true},
		{name: "invalid source", envelope: &meshtastic.ServiceEnvelope{
			ChannelId: "channel",
			GatewayId: "gateway",
			Packet: &meshtastic.MeshPacket{
				Id:        1,
				Encrypted: []byte{1},
			},
		}, wantErr: true},
		{name: "missing encrypted payload", envelope: &meshtastic.ServiceEnvelope{
			ChannelId: "channel",
			GatewayId: "gateway",
			Packet: &meshtastic.MeshPacket{
				Id:   1,
				From: 2,
			},
		}, wantErr: true},
		{name: "already decoded", envelope: &meshtastic.ServiceEnvelope{
			ChannelId: "channel",
			GatewayId: "gateway",
			Packet: &meshtastic.MeshPacket{
				Id:        1,
				From:      2,
				Encrypted: []byte{1},
				Decoded:   &meshtastic.Data{Payload: []byte("decoded")},
			},
		}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateServiceEnvelope(tt.envelope); (err != nil) != tt.wantErr {
				t.Fatalf("validateServiceEnvelope() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecryptMeshPacket(t *testing.T) {
	handler := testHandler()
	want := &meshtastic.Data{
		Portnum: 1,
		Payload: []byte("hello meshtastic"),
	}

	envelope := testEnvelope(t, want)
	got, err := handler.decryptMeshPacket(envelope)
	if err != nil {
		t.Fatalf("decryptMeshPacket() retornou erro: %v", err)
	}

	if !proto.Equal(got, want) {
		t.Fatalf("payload decodificado inesperado: got %v want %v", got, want)
	}
}

func TestDecryptMeshPacketReturnsErrorForInvalidCiphertext(t *testing.T) {
	handler := testHandler()
	envelope := &meshtastic.ServiceEnvelope{
		ChannelId: "channel",
		GatewayId: "gateway",
		Packet: &meshtastic.MeshPacket{
			Id:        1,
			From:      2,
			Encrypted: []byte{1, 2, 3, 4},
		},
	}

	got, err := handler.decryptMeshPacket(envelope)
	if err == nil {
		t.Fatal("decryptMeshPacket() deveria retornar erro para ciphertext inválido")
	}
	if got != nil {
		t.Fatalf("decryptMeshPacket() retornou dados inesperados: %v", got)
	}
}

func TestDecryptMeshPacketReturnsNilForUnknownData(t *testing.T) {
	handler := testHandler()
	envelope := testEnvelope(t, &meshtastic.Data{Portnum: 0, Payload: []byte("unknown")})

	got, err := handler.decryptMeshPacket(envelope)
	if err != nil {
		t.Fatalf("decryptMeshPacket() deveria aceitar Data não reconhecido: %v", err)
	}
	if got != nil {
		t.Fatalf("decryptMeshPacket() deveria retornar nil para Data não reconhecido: %v", got)
	}
}

func TestOnPublishIgnoresEmptyPayload(t *testing.T) {
	handler := testHandler()
	client := &mqtt.Client{ID: "test-client"}

	got, err := handler.OnPublish(client, packets.Packet{
		TopicName: "msh/test",
	})

	if err != packets.CodeSuccessIgnore {
		t.Fatalf("OnPublish() erro = %v, want %v", err, packets.CodeSuccessIgnore)
	}
	if got.Payload != nil {
		t.Fatalf("OnPublish() alterou payload: %x", got.Payload)
	}
}

func TestOnPublishIgnoresInvalidProtobuf(t *testing.T) {
	handler := testHandler()
	client := &mqtt.Client{ID: "test-client"}
	payload := []byte{0xff, 0xff, 0xff}

	got, err := handler.OnPublish(client, packets.Packet{
		TopicName: "msh/test",
		Payload:   payload,
	})

	if err != packets.CodeSuccessIgnore {
		t.Fatalf("OnPublish() erro = %v, want %v", err, packets.CodeSuccessIgnore)
	}
	if !bytes.Equal(got.Payload, payload) {
		t.Fatalf("OnPublish() alterou payload: %x != %x", got.Payload, payload)
	}
}

func TestOnPublishAcceptsValidEnvelope(t *testing.T) {
	handler := testHandler()
	client := &mqtt.Client{ID: "test-client"}

	envelope := testEnvelope(t, &meshtastic.Data{
		Portnum: 1,
		Payload: []byte("hello"),
	})
	payload, err := proto.Marshal(envelope)
	if err != nil {
		t.Fatalf("proto.Marshal() retornou erro: %v", err)
	}

	got, err := handler.OnPublish(client, packets.Packet{
		TopicName: "msh/test",
		Payload:   payload,
	})

	if err != nil {
		t.Fatalf("OnPublish() retornou erro: %v", err)
	}
	if !bytes.Equal(got.Payload, payload) {
		t.Fatalf("OnPublish() alterou payload")
	}
}

func TestOnPublishBlocksEnvelopeWhenDecryptionCannotDecodeData(t *testing.T) {
	handler := testHandler()
	client := &mqtt.Client{ID: "test-client"}

	envelope := &meshtastic.ServiceEnvelope{
		ChannelId: "channel",
		GatewayId: "gateway",
		Packet: &meshtastic.MeshPacket{
			Id:        1,
			From:      2,
			Encrypted: []byte{1, 2, 3, 4},
		},
	}
	payload, err := proto.Marshal(envelope)
	if err != nil {
		t.Fatalf("proto.Marshal() retornou erro: %v", err)
	}

	got, err := handler.OnPublish(client, packets.Packet{
		TopicName: "msh/test",
		Payload:   payload,
	})

	if err != packets.CodeSuccessIgnore {
		t.Fatalf("OnPublish() deveria bloquear falha de parsing do Data: got %v", err)
	}
	if !bytes.Equal(got.Payload, payload) {
		t.Fatalf("OnPublish() alterou payload")
	}
}
