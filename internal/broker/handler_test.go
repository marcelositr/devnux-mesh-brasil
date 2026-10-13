package broker

import (
	"bytes"
	"io"
	"log/slog"
	"strings"
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


func TestLogReceivedMessageDoesNotLogTextPayload(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, nil))
	handler := NewHandler(logger, config.Config{DefaultPSK: crypto.DefaultPSK})
	secretMessage := "texto privado que nao deve aparecer no log"

	handler.logReceivedMessage("msh/test", "test-client", &meshtastic.Data{
		Portnum: 1,
		Payload: []byte(secretMessage),
	})

	output := logs.String()
	if strings.Contains(output, secretMessage) {
		t.Fatalf("o log revelou o conteúdo da mensagem: %q", output)
	}
	if !strings.Contains(output, "payload_bytes=") {
		t.Fatalf("o log deveria registrar o tamanho do payload: %q", output)
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
			Id:   1,
			From: 2,
			PayloadVariant: &meshtastic.MeshPacket_Encrypted{
				Encrypted: []byte{1, 2, 3, 4},
			},
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
