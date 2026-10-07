package crypto

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
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

func TestTransformPacketMatchesInteroperabilityVector(t *testing.T) {
	// O vetor valida a interoperabilidade da cifra com uma implementação C#.
	// Os valores de origem e identificador correspondem ao pacote usado na validação.
	nonce := NewNonce(4202784164, 1777428186)

	ciphertext, err := base64.StdEncoding.DecodeString("kiDV39nDDsi8AON+Czei6zUpy+F/7E+lyIpicxJR40KXBFmPkqFUEnobI5voQadha+s=")
	if err != nil {
		t.Fatalf("ciphertext C# inválido: %v", err)
	}

	expected, err := hex.DecodeString("0804122c0a09216661383136356134120f4d65736874617374696320363561341a04363561342206f412fa8165a4282b1801")
	if err != nil {
		t.Fatalf("plaintext esperado inválido: %v", err)
	}

	decrypted, err := TransformPacket(ciphertext, nonce, DefaultPSK)
	if err != nil {
		t.Fatalf("erro ao descriptografar vetor de interoperabilidade: %v", err)
	}

	if !bytes.Equal(decrypted, expected) {
		t.Fatalf("vetor de interoperabilidade incompatível: %x != %x", decrypted, expected)
	}
}
