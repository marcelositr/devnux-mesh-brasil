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

func TestTransformPacketMatchesMeshtasticCSharpVector(t *testing.T) {
	// Vetor retirado de Meshtastic/c-sharp:
	// Meshtastic.Test/Crypto/PacketCryptoTests.cs
	// packet.From = 4202784164
	// packet.Id = 1777428186
	// O ciphertext foi gerado pela implementação C# usando Resources.DEFAULT_PSK.
	nonce := NewNonce(4202784164, 1777428186)

	ciphertext := []byte{
		0x92, 0x20, 0xd5, 0xdf, 0xd9, 0xc3, 0x0e, 0xc8,
		0x7c, 0x00, 0xe3, 0x7f, 0x0b, 0x27, 0xa2, 0xeb,
		0x35, 0x29, 0x31, 0x7b, 0xe3, 0x7f, 0xec, 0x4f,
		0xa5, 0xc8, 0x98, 0x9c, 0x4f, 0x47, 0x4a, 0x50,
		0x17, 0x63, 0xe3, 0xa1, 0x30, 0xa6, 0x1a, 0x75,
		0xeb, 0x1c,
	}

	decrypted, err := TransformPacket(ciphertext, nonce, DefaultPSK)
	if err != nil {
		t.Fatalf("erro ao descriptografar vetor C#: %v", err)
	}

	expected := []byte{
		0x08, 0x04, 0x12, 0x2c,
		0x0a, 0x09, 0x21, 0x66, 0x61, 0x38, 0x31, 0x36, 0x35, 0x61, 0x34,
		0x12, 0x0f, 0x4d, 0x65, 0x73, 0x68, 0x74, 0x61, 0x73, 0x74, 0x69, 0x63, 0x20, 0x36, 0x35, 0x61, 0x34,
		0x1a, 0x04, 0x36, 0x35, 0x61, 0x34,
		0x22, 0x06, 0xf4, 0x12, 0xfa, 0x81, 0x65, 0xa4,
		0x28, 0x2b, 0x18, 0x01,
	}

	if !bytes.Equal(decrypted, expected) {
		t.Fatalf("vetor C# incompatível: %x != %x", decrypted, expected)
	}
}
