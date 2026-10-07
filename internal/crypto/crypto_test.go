package crypto

import (
	"bytes"
	"testing"
)

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
