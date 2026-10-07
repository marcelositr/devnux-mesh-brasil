package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
)

func NewNonce(from, packetID uint32) []byte {
	nonce := make([]byte, aes.BlockSize)
	binary.LittleEndian.PutUint32(nonce[0:4], packetID)
	binary.LittleEndian.PutUint32(nonce[4:8], from)
	return nonce
}

func TransformPacket(input, nonce, psk []byte) ([]byte, error) {
	if len(nonce) != aes.BlockSize {
		return nil, fmt.Errorf("nonce inválido: esperado %d bytes, recebido %d", aes.BlockSize, len(nonce))
	}

	key, err := normalizeKey(psk)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("criar AES: %w", err)
	}

	output := make([]byte, len(input))
	stream := cipher.NewCTR(block, nonce)
	stream.XORKeyStream(output, input)
	return output, nil
}

func normalizeKey(psk []byte) ([]byte, error) {
	switch len(psk) {
	case 16, 32:
		key := make([]byte, len(psk))
		copy(key, psk)
		return key, nil
	case 0:
		return nil, fmt.Errorf("PSK vazia")
	default:
		if len(psk) < 16 {
			key := make([]byte, 16)
			copy(key, psk)
			return key, nil
		}
		return nil, fmt.Errorf("tamanho de PSK não suportado: %d", len(psk))
	}
}
