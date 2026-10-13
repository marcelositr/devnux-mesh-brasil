package broker

import (
	"bytes"
	"io"
	"log/slog"
	"testing"

	meshtastic "github.com/kmpm/meshtastic-protobufs.go/v2/generated"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
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

func marshalEnvelope(t *testing.T, envelope *meshtastic.ServiceEnvelope) []byte {
	t.Helper()
	payload, err := proto.Marshal(envelope)
	if err != nil {
		t.Fatalf("proto.Marshal() retornou erro: %v", err)
	}
	return payload
}

func TestValidateServiceEnvelope(t *testing.T) {
	processor := newPacketProcessor(crypto.DefaultPSK)
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
		name     string
		envelope *meshtastic.ServiceEnvelope
		wantErr  bool
	}{
		{name: "valid", envelope: valid},
		{name: "nil", envelope: nil, wantErr: true},
		{name: "missing channel", envelope: &meshtastic.ServiceEnvelope{
			GatewayId: "gateway", Packet: valid.Packet,
		}, wantErr: true},
		{name: "missing gateway", envelope: &meshtastic.ServiceEnvelope{
			ChannelId: "channel", Packet: valid.Packet,
		}, wantErr: true},
		{name: "missing packet", envelope: &meshtastic.ServiceEnvelope{
			ChannelId: "channel", GatewayId: "gateway",
		}, wantErr: true},
		{name: "invalid packet id", envelope: &meshtastic.ServiceEnvelope{
			ChannelId: "channel", GatewayId: "gateway",
			Packet: &meshtastic.MeshPacket{From: 2, Encrypted: []byte{1}},
		}, wantErr: true},
		{name: "invalid source", envelope: &meshtastic.ServiceEnvelope{
			ChannelId: "channel", GatewayId: "gateway",
			Packet: &meshtastic.MeshPacket{Id: 1, Encrypted: []byte{1}},
		}, wantErr: true},
		{name: "missing payload", envelope: &meshtastic.ServiceEnvelope{
			ChannelId: "channel", GatewayId: "gateway",
			Packet: &meshtastic.MeshPacket{Id: 1, From: 2},
		}, wantErr: true},
		{name: "decoded packet accepted for inspection", envelope: &meshtastic.ServiceEnvelope{
			ChannelId: "channel", GatewayId: "gateway",
			Packet: &meshtastic.MeshPacket{Id: 1, From: 2, Decoded: &meshtastic.Data{Payload: []byte("decoded")}},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := processor.validateServiceEnvelope(tt.envelope); (err != nil) != tt.wantErr {
				t.Fatalf("validateServiceEnvelope() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecryptMeshPacket(t *testing.T) {
	processor := newPacketProcessor(crypto.DefaultPSK)
	want := &meshtastic.Data{Portnum: 1, Payload: []byte("hello meshtastic")}
	got, err := processor.decryptMeshPacket(testEnvelope(t, want))
	if err != nil {
		t.Fatalf("decryptMeshPacket() retornou erro: %v", err)
	}
	if !proto.Equal(got, want) {
		t.Fatalf("payload decodificado inesperado: got %v want %v", got, want)
	}
}

func TestDecryptMeshPacketReturnsErrorForInvalidCiphertext(t *testing.T) {
	processor := newPacketProcessor(crypto.DefaultPSK)
	envelope := &meshtastic.ServiceEnvelope{
		ChannelId: "channel", GatewayId: "gateway",
		Packet: &meshtastic.MeshPacket{Id: 1, From: 2, Encrypted: []byte{1, 2, 3, 4}},
	}
	if got, err := processor.decryptMeshPacket(envelope); err == nil || got != nil {
		t.Fatalf("decryptMeshPacket() = (%v, %v), esperava erro e dados nil", got, err)
	}
}

func TestDecryptMeshPacketReturnsNilForUnknownData(t *testing.T) {
	processor := newPacketProcessor(crypto.DefaultPSK)
	envelope := testEnvelope(t, &meshtastic.Data{Portnum: 0, Payload: []byte("unknown")})
	if got, err := processor.decryptMeshPacket(envelope); err != nil || got != nil {
		t.Fatalf("decryptMeshPacket() = (%v, %v), esperava (nil, nil)", got, err)
	}
}

func TestOnPublishForwardsJSONWithoutProtobufParsing(t *testing.T) {
	handler := testHandler()
	client := &mqtt.Client{ID: "test-client"}
	payload := []byte(`{"from":2740626832,"type":"sendtext","payload":"TESTE TERMINAL"}`)

	got, err := handler.OnPublish(client, packets.Packet{
		TopicName: "msh/2/json/mqtt/",
		Payload:   payload,
	})
	if err != nil {
		t.Fatalf("OnPublish() retornou erro: %v", err)
	}
	if !bytes.Equal(got.Payload, payload) {
		t.Fatalf("payload JSON foi alterado: %q != %q", got.Payload, payload)
	}
}

func TestOnPublishPreservesInvalidProtobufOnMeshtasticTopic(t *testing.T) {
	handler := testHandler()
	client := &mqtt.Client{ID: "test-client"}
	payload := []byte{0xff, 0xff, 0xff}

	got, err := handler.OnPublish(client, packets.Packet{
		TopicName: "msh/2/e/ITUVERAVA/!abc",
		Payload:   payload,
	})
	if err != nil {
		t.Fatalf("OnPublish() deveria preservar a publicação: %v", err)
	}
	if !bytes.Equal(got.Payload, payload) {
		t.Fatalf("payload foi alterado: %x != %x", got.Payload, payload)
	}
}

func TestOnPublishPreservesPacketWhenLocalDecryptionFails(t *testing.T) {
	handler := testHandler()
	client := &mqtt.Client{ID: "test-client"}
	envelope := &meshtastic.ServiceEnvelope{
		ChannelId: "channel", GatewayId: "gateway",
		Packet: &meshtastic.MeshPacket{Id: 1, From: 2, Encrypted: []byte{1, 2, 3, 4}},
	}
	payload := marshalEnvelope(t, envelope)

	got, err := handler.OnPublish(client, packets.Packet{
		TopicName: "msh/2/e/ITUVERAVA/!abc",
		Payload:   payload,
	})
	if err != nil {
		t.Fatalf("OnPublish() não deve bloquear falha de inspeção: %v", err)
	}
	if !bytes.Equal(got.Payload, payload) {
		t.Fatal("payload original não foi preservado")
	}
}

func TestOnPublishZeroHopRewritesOnlyHopLimit(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewHandler(logger, config.Config{DefaultPSK: crypto.DefaultPSK, ZeroHop: true})
	client := &mqtt.Client{ID: "test-client"}

	envelope := testEnvelope(t, &meshtastic.Data{Portnum: 1, Payload: []byte("hello")})
	envelope.GetPacket().HopLimit = 3
	envelope.GetPacket().HopStart = 4
	payload := marshalEnvelope(t, envelope)

	got, err := handler.OnPublish(client, packets.Packet{
		TopicName: "msh/2/e/ITUVERAVA/!abc",
		Payload:   payload,
	})
	if err != nil {
		t.Fatalf("OnPublish() retornou erro: %v", err)
	}

	var rewritten meshtastic.ServiceEnvelope
	if err := proto.Unmarshal(got.Payload, &rewritten); err != nil {
		t.Fatalf("payload reescrito não é protobuf válido: %v", err)
	}
	if rewritten.GetPacket().GetHopLimit() != 0 {
		t.Fatalf("hop_limit = %d, esperado 0", rewritten.GetPacket().GetHopLimit())
	}
	if rewritten.GetPacket().GetHopStart() != 4 {
		t.Fatalf("hop_start = %d, esperado 4", rewritten.GetPacket().GetHopStart())
	}
}

func TestIsProtobufTopic(t *testing.T) {
	tests := []struct {
		topic string
		want  bool
	}{
		{"msh/2/e/ITUVERAVA/!a35aa590", true},
		{"msh/2/c/LongFast/!a35aa590", true},
		{"msh/2/json/mqtt/", false},
		{"msh/2/json/LongFast/!abc", false},
		{"other/topic", false},
	}
	for _, tt := range tests {
		t.Run(tt.topic, func(t *testing.T) {
			if got := isProtobufTopic(tt.topic); got != tt.want {
				t.Fatalf("isProtobufTopic(%q) = %v, esperado %v", tt.topic, got, tt.want)
			}
		})
	}
}
