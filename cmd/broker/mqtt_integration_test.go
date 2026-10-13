package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"encoding/pem"
	"io"
	"log/slog"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	meshtastic "github.com/kmpm/meshtastic-protobufs.go/v2/generated"
	"github.com/marcelositr/devnux-mesh-brasil/internal/config"
	"github.com/marcelositr/devnux-mesh-brasil/internal/crypto"
	"google.golang.org/protobuf/proto"
)

func TestMQTTSubscriberReceivesAcceptedPublication(t *testing.T) {
	certificateFile, privateKeyFile, certificate := writeTestCertificate(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server, err := newServer(logger, config.Config{
		ListenAddress:   "127.0.0.1:0",
		CertificateFile: certificateFile,
		PrivateKeyFile:  privateKeyFile,
		DefaultPSK:      crypto.DefaultPSK,
	})
	if err != nil {
		t.Fatalf("newServer() retornou erro: %v", err)
	}
	t.Cleanup(func() {
		if err := server.Close(); err != nil {
			t.Errorf("server.Close() retornou erro: %v", err)
		}
	})

	listener, ok := server.Listeners.Get("mqtt-tls")
	if !ok {
		t.Fatal("listener mqtt-tls não foi registrado")
	}

	if err := server.Serve(); err != nil {
		t.Fatalf("server.Serve() retornou erro: %v", err)
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    certificatePool(t, certificate),
		ServerName: "127.0.0.1",
	}

	subscriber := connectTestMQTTClient(t, listener.Address(), tlsConfig, "integration-subscriber")
	t.Cleanup(func() { _ = subscriber.Close() })
	subscribeTestMQTTClient(t, subscriber, "msh/test")

	publisher := connectTestMQTTClient(t, listener.Address(), tlsConfig, "integration-publisher")
	t.Cleanup(func() { _ = publisher.Close() })

	const topic = "msh/test"
	payload := validTestEnvelopePayload(t)
	if err := writeMQTTPacket(publisher, 0x30, encodeMQTTString(topic, payload)); err != nil {
		t.Fatalf("não foi possível publicar mensagem MQTT: %v", err)
	}

	if err := subscriber.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatalf("não foi possível definir timeout: %v", err)
	}
	header, body, err := readMQTTPacket(subscriber)
	if err != nil {
		t.Fatalf("assinante não recebeu publicação MQTT: %v", err)
	}
	if header>>4 != 3 {
		t.Fatalf("pacote recebido não é PUBLISH: header=0x%02x", header)
	}

	receivedTopic, receivedPayload, err := decodeMQTTPublish(body)
	if err != nil {
		t.Fatalf("não foi possível decodificar PUBLISH recebido: %v", err)
	}
	if receivedTopic != topic {
		t.Fatalf("tópico recebido = %q, esperado %q", receivedTopic, topic)
	}
	if !bytes.Equal(receivedPayload, payload) {
		t.Fatal("payload recebido pelo assinante difere do payload publicado")
	}
}

func connectTestMQTTClient(t *testing.T, address string, tlsConfig *tls.Config, clientID string) net.Conn {
	t.Helper()

	conn, err := tls.Dial("tcp", address, tlsConfig.Clone())
	if err != nil {
		t.Fatalf("tls.Dial() retornou erro: %v", err)
	}

	body := append(encodeMQTTString("MQTT"), 4, 2, 0, 30)
	body = append(body, encodeMQTTString(clientID)...)
	if err := writeMQTTPacket(conn, 0x10, body); err != nil {
		_ = conn.Close()
		t.Fatalf("não foi possível enviar CONNECT: %v", err)
	}

	header, response, err := readMQTTPacket(conn)
	if err != nil {
		_ = conn.Close()
		t.Fatalf("não foi possível ler CONNACK: %v", err)
	}
	if header>>4 != 2 || len(response) != 2 || response[1] != 0 {
		_ = conn.Close()
		t.Fatalf("CONNACK inesperado: header=0x%02x body=%x", header, response)
	}

	return conn
}

func subscribeTestMQTTClient(t *testing.T, conn net.Conn, topic string) {
	t.Helper()

	body := []byte{0, 1}
	body = append(body, encodeMQTTString(topic)...)
	body = append(body, 0)
	if err := writeMQTTPacket(conn, 0x82, body); err != nil {
		t.Fatalf("não foi possível enviar SUBSCRIBE: %v", err)
	}

	if err := conn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatalf("não foi possível definir timeout de SUBACK: %v", err)
	}
	header, response, err := readMQTTPacket(conn)
	if err != nil {
		t.Fatalf("não foi possível ler SUBACK: %v", err)
	}
	if header>>4 != 9 || len(response) < 3 || response[2] > 1 {
		t.Fatalf("SUBACK inesperado: header=0x%02x body=%x", header, response)
	}
	if err := conn.SetReadDeadline(time.Time{}); err != nil {
		t.Fatalf("não foi possível remover timeout: %v", err)
	}
}

func writeMQTTPacket(conn net.Conn, header byte, body []byte) error {
	var packet bytes.Buffer
	packet.WriteByte(header)
	remaining := len(body)
	for {
		encoded := byte(remaining % 128)
		remaining /= 128
		if remaining > 0 {
			encoded |= 0x80
		}
		packet.WriteByte(encoded)
		if remaining == 0 {
			break
		}
	}
	packet.Write(body)
	_, err := conn.Write(packet.Bytes())
	return err
}

func readMQTTPacket(conn net.Conn) (byte, []byte, error) {
	var header [1]byte
	if _, err := io.ReadFull(conn, header[:]); err != nil {
		return 0, nil, err
	}

	remaining := 0
	multiplier := 1
	for i := 0; i < 4; i++ {
		var encoded [1]byte
		if _, err := io.ReadFull(conn, encoded[:]); err != nil {
			return 0, nil, err
		}
		remaining += int(encoded[0]&127) * multiplier
		if encoded[0]&128 == 0 {
			body := make([]byte, remaining)
			if _, err := io.ReadFull(conn, body); err != nil {
				return 0, nil, err
			}
			return header[0], body, nil
		}
		multiplier *= 128
	}
	return 0, nil, os.ErrInvalid
}

func encodeMQTTString(value string, suffix ...[]byte) []byte {
	encoded := make([]byte, 2+len(value))
	binary.BigEndian.PutUint16(encoded, uint16(len(value)))
	copy(encoded[2:], value)
	for _, part := range suffix {
		encoded = append(encoded, part...)
	}
	return encoded
}

func decodeMQTTPublish(body []byte) (string, []byte, error) {
	if len(body) < 2 {
		return "", nil, io.ErrUnexpectedEOF
	}
	topicLength := int(binary.BigEndian.Uint16(body[:2]))
	if len(body) < 2+topicLength {
		return "", nil, io.ErrUnexpectedEOF
	}
	return string(body[2 : 2+topicLength]), body[2+topicLength:], nil
}

func validTestEnvelopePayload(t *testing.T) []byte {
	t.Helper()

	const (
		from     = uint32(12345)
		packetID = uint32(67890)
	)
	plain, err := proto.Marshal(&meshtastic.Data{
		Portnum: 1,
		Payload: []byte("integration-test-message"),
	})
	if err != nil {
		t.Fatalf("não foi possível serializar Data: %v", err)
	}
	encrypted, err := crypto.TransformPacket(plain, crypto.NewNonce(from, packetID), crypto.DefaultPSK)
	if err != nil {
		t.Fatalf("não foi possível criptografar Data: %v", err)
	}
	envelope, err := proto.Marshal(&meshtastic.ServiceEnvelope{
		ChannelId: "integration-channel",
		GatewayId: "integration-gateway",
		Packet: &meshtastic.MeshPacket{
			Id:   packetID,
			From: from,
			PayloadVariant: &meshtastic.MeshPacket_Encrypted{
				Encrypted: encrypted,
			},
		},
	})
	if err != nil {
		t.Fatalf("não foi possível serializar ServiceEnvelope: %v", err)
	}
	return envelope
}

func writeTestCertificate(t *testing.T) (string, string, *x509.Certificate) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("não foi possível gerar chave de teste: %v", err)
	}
	now := time.Now()
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "127.0.0.1"},
		NotBefore:             now.Add(-time.Minute),
		NotAfter:              now.Add(time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("não foi possível criar certificado de teste: %v", err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("não foi possível interpretar certificado de teste: %v", err)
	}
	privateKey, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("não foi possível serializar chave de teste: %v", err)
	}

	certificateFile := filepath.Join(t.TempDir(), "certificate.pem")
	privateKeyFile := filepath.Join(t.TempDir(), "private.key")
	if err := os.WriteFile(certificateFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0600); err != nil {
		t.Fatalf("não foi possível gravar certificado de teste: %v", err)
	}
	if err := os.WriteFile(privateKeyFile, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKey}), 0600); err != nil {
		t.Fatalf("não foi possível gravar chave de teste: %v", err)
	}
	return certificateFile, privateKeyFile, certificate
}

func certificatePool(t *testing.T, certificate *x509.Certificate) *x509.CertPool {
	t.Helper()
	pool := x509.NewCertPool()
	pool.AddCert(certificate)
	return pool
}
