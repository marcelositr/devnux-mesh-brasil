package crypto

import (
	"bytes"
	"testing"
)

func TestNewNonceMatchesMeshtasticLayout(t *testing.T) {
	nonce := NewNonce(0x0a0b0c0d, 0x01020304)

	expected := []byte{
		0x04, 0x03, 0x02, 0x01,
		0x00, 0x00, 0x00, 0x00,
		0x0d, 0x0c, 0x0b, 0x0a,
		0x00, 0x00, 0x00, 0x00,
	}

	if !bytes.Equal(nonce, expected) {
		t.Fatalf("nonce incompatível com Meshtastic: %x != %x", nonce, expected)
	}
}

func TestDefaultPSK(t *testing.T) {
	expected := []byte{
		0xd4, 0xf1, 0xbb, 0x3a,
		0x20, 0x29, 0x07, 0x59,
		0xf0, 0xbc, 0xff, 0xab,
		0xcf, 0x4e, 0x69, 0x01,
	}

	if !bytes.Equal(DefaultPSK, expected) {
		t.Fatalf("chave padrão incompatível com Meshtastic: %x != %x", DefaultPSK, expected)
	}
}

func TestTransformPacketRoundTrip(t *testing.T) {
	nonce := NewNonce(0x01020304, 0x05060708)
	key := []byte("0123456789abcdef")
	plain := []byte("meshtastic")

	encrypted, err := TransformPacket(plain, nonce, key)
	if err != nil {
		t.Fatalf("erro ao criptografar: %v", err)
	}

	decrypted, err := TransformPacket(encrypted, nonce, key)
	if err != nil {
		t.Fatalf("erro ao descriptografar: %v", err)
	}

	if !bytes.Equal(plain, decrypted) {
		t.Fatalf("payload diferente após round trip: %q != %q", plain, decrypted)
	}
}

func TestTransformPacketRejectsUnsupportedKeySize(t *testing.T) {
	nonce := NewNonce(0x01020304, 0x05060708)

	if _, err := TransformPacket([]byte("meshtastic"), nonce, []byte("chave-curta")); err == nil {
		t.Fatal("esperava erro para chave com tamanho não suportado")
	}
}
