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

func TestOnPublishForwardsEmptyPayload(t *testing.T) {
	handler := testHandler()
	client := &mqtt.Client{ID: "test-client"}

	got, err := handler.OnPublish(client, packets.Packet{
		TopicName: "msh/2/e/LongFast/!abcd1234",
	})

	if err != nil {
		t.Fatalf("OnPublish() deveria encaminhar payload sem erro, recebeu %v", err)
	}
	if got.Payload != nil {
		t.Fatalf("OnPublish() alterou payload: %x", got.Payload)
	}
}

func TestOnPublishForwardsInvalidProtobuf(t *testing.T) {
	handler := testHandler()
	client := &mqtt.Client{ID: "test-client"}
	payload := []byte{0xff, 0xff, 0xff}

	got, err := handler.OnPublish(client, packets.Packet{
		TopicName: "msh/2/e/LongFast/!abcd1234",
		Payload:   payload,
	})

	if err != nil {
		t.Fatalf("OnPublish() deveria encaminhar protobuf inválido sem erro, recebeu %v", err)
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
		TopicName: "msh/2/e/LongFast/!abcd1234",
		Payload:   payload,
	})

	if err != nil {
		t.Fatalf("OnPublish() retornou erro: %v", err)
	}
	if !bytes.Equal(got.Payload, payload) {
		t.Fatalf("OnPublish() alterou payload")
	}
}

func TestOnPublishForwardsEnvelopeWhenDecryptionCannotDecodeData(t *testing.T) {
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
		TopicName: "msh/2/e/LongFast/!abcd1234",
		Payload:   payload,
	})

	if err != nil {
		t.Fatalf("OnPublish() deveria preservar o envelope quando a descriptografia falha: %v", err)
	}
	if !bytes.Equal(got.Payload, payload) {
		t.Fatalf("OnPublish() alterou payload")
	}
}

func TestOnPublishForwardsDataWithoutUsefulPayload(t *testing.T) {
	tests := []struct {
		name string
		data *meshtastic.Data
	}{
		{
			name: "zero portnum",
			data: &meshtastic.Data{Portnum: 0, Payload: []byte("non-empty")},
		},
		{
			name: "empty payload",
			data: &meshtastic.Data{Portnum: 1},
		},
		{
			name: "zero portnum and empty payload",
			data: &meshtastic.Data{Portnum: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := testHandler()
			client := &mqtt.Client{ID: "test-client"}
			payload, err := proto.Marshal(testEnvelope(t, tt.data))
			if err != nil {
				t.Fatalf("proto.Marshal() retornou erro: %v", err)
			}

			got, result := handler.OnPublish(client, packets.Packet{
				TopicName: "msh/2/e/LongFast/!abcd1234",
				Payload:   payload,
			})
			if result != nil {
				t.Fatalf("OnPublish() deveria encaminhar Data sem payload útil, recebeu %v", result)
			}
			if !bytes.Equal(got.Payload, payload) {
				t.Fatal("OnPublish() alterou o payload original")
			}
		})
	}
}

func TestOnPublishAcceptsUnknownPortnumWithPayload(t *testing.T) {
	handler := testHandler()
	client := &mqtt.Client{ID: "test-client"}
	envelope := testEnvelope(t, &meshtastic.Data{
		Portnum: meshtastic.PortNum(999),
		Payload: []byte("valid unknown portnum payload"),
	})
	payload, err := proto.Marshal(envelope)
	if err != nil {
		t.Fatalf("proto.Marshal() retornou erro: %v", err)
	}

	got, result := handler.OnPublish(client, packets.Packet{
		TopicName: "msh/2/e/LongFast/!abcd1234",
		Payload:   payload,
	})
	if result != nil {
		t.Fatalf("OnPublish() deveria aceitar portnum desconhecido com payload: %v", result)
	}
	if !bytes.Equal(got.Payload, payload) {
		t.Fatal("OnPublish() alterou o payload original")
	}
}


func TestOnPublishForwardsJSONDownlinkWithoutParsingAsServiceEnvelope(t *testing.T) {
	handler := testHandler()
	client := &mqtt.Client{ID: "test-client"}
	payload := []byte(`{"from":12345,"type":"sendtext","payload":"integration test"}`)
	topic := "msh/2/json/mqtt/"

	got, err := handler.OnPublish(client, packets.Packet{
		TopicName: topic,
		Payload:   payload,
	})
	if err != nil {
		t.Fatalf("OnPublish() deveria preservar tópico JSON: %v", err)
	}
	if got.TopicName != topic || !bytes.Equal(got.Payload, payload) {
		t.Fatal("OnPublish() alterou o tópico ou payload JSON")
	}
}

func TestIsProtobufTopic(t *testing.T) {
	tests := []struct {
		topic string
		want  bool
	}{
		{topic: "msh/BR/2/e/LongFast/!abcd1234", want: true},
		{topic: "msh/2/e/LongFast/!abcd1234", want: true},
		{topic: "msh/US/2/c/LongFast/!abcd1234", want: true},
		{topic: "msh/BR/2/json/LongFast/!abcd1234", want: false},
		{topic: "msh/BR/2/map/!abcd1234", want: false},
		{topic: "msh/BR/2/stat/!abcd1234", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.topic, func(t *testing.T) {
			if got := isProtobufTopic(tt.topic); got != tt.want {
				t.Fatalf("isProtobufTopic(%q) = %v, want %v", tt.topic, got, tt.want)
			}
		})
	}
}
