package config

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/marcelositr/devnux-mesh-brasil/internal/crypto"
)

// Config reúne as configurações usadas para iniciar o broker.
type Config struct {
	ListenAddress   string
	CertificateFile string
	PrivateKeyFile  string
	DefaultPSK      []byte
}

// Load carrega as configurações do ambiente e aplica os valores padrão.
func Load() (Config, error) {
	cfg := Config{
		ListenAddress:   envOrDefault("MQTT_LISTEN_ADDRESS", ":8883"),
		CertificateFile: envOrDefault("MQTT_CERTIFICATE_FILE", "certificate.pem"),
		PrivateKeyFile:  envOrDefault("MQTT_PRIVATE_KEY_FILE", "private.key"),
		DefaultPSK:      append([]byte(nil), crypto.DefaultPSK...),
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

// parsePSK converte a PSK hexadecimal e valida os tamanhos aceitos pelo AES.
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

// envOrDefault retorna o valor da variável de ambiente ou o fallback quando ela não está definida.
func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
