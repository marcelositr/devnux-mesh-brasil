package config

import (
	"encoding/hex"
	"fmt"
	"os"
)

type Config struct {
	ListenAddress   string
	CertificateFile string
	PrivateKeyFile  string
	DefaultPSK      []byte
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

func parsePSK(value string) ([]byte, error) {
	psk, err := hex.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("esperado hexadecimal: %w", err)
	}

	switch len(psk) {
	case 16, 32:
		return psk, nil
	default:
		return nil, fmt.Errorf("tamanho inválido: esperado 16 ou 32 bytes, recebido %d", len(psk))
	}
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
