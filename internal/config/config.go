package config

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ListenAddress   string
	CertificateFile string
	PrivateKeyFile  string
	DefaultPSK      []byte
	ZeroHop         bool
}

func Load() (Config, error) {
	cfg := Config{
		ListenAddress:   envOrDefault("MQTT_LISTEN_ADDRESS", ":8883"),
		CertificateFile: envOrDefault("MQTT_CERTIFICATE_FILE", "certificate.pem"),
		PrivateKeyFile:  envOrDefault("MQTT_PRIVATE_KEY_FILE", "private.key"),
		DefaultPSK:      defaultPSK(),
	}

	if value := os.Getenv("MESHTASTIC_DEFAULT_PSK"); value != "" {
		psk, err := parsePSK(value)
		if err != nil {
			return Config{}, fmt.Errorf("MESHTASTIC_DEFAULT_PSK inválida: %w", err)
		}
		cfg.DefaultPSK = psk
	}

	if value := os.Getenv("MESHTASTIC_ZERO_HOP"); value != "" {
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return Config{}, fmt.Errorf("MESHTASTIC_ZERO_HOP inválido: %w", err)
		}
		cfg.ZeroHop = enabled
	}

	return cfg, nil
}

func defaultPSK() []byte {
	return []byte{
		0xd4, 0xf1, 0xbb, 0x3a,
		0x20, 0x29, 0x07, 0x59,
		0xf0, 0xbc, 0xff, 0xab,
		0xcf, 0x4e, 0x69, 0x01,
	}
}

// parsePSK accepts the hexadecimal format used by previous versions and
// standard Base64, as commonly displayed by Meshtastic for 256-bit keys.
func parsePSK(value string) ([]byte, error) {
	if psk, err := hex.DecodeString(value); err == nil {
		if validPSKSize(psk) {
			return psk, nil
		}
	}

	psk, base64Err := base64.StdEncoding.DecodeString(value)
	if base64Err == nil && validPSKSize(psk) {
		return psk, nil
	}

	if base64Err != nil {
		return nil, fmt.Errorf("esperado hexadecimal ou Base64 válido de 16 ou 32 bytes")
	}
	return nil, fmt.Errorf("tamanho inválido: esperado 16 ou 32 bytes, recebido %d", len(psk))
}

func validPSKSize(psk []byte) bool {
	return len(psk) == 16 || len(psk) == 32
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
