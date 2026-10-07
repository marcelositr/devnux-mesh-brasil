package config

import (
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
		DefaultPSK:      []byte(envOrDefault("MESHTASTIC_DEFAULT_PSK", "meshtastic_ota_default_psk_v1!!!")),
	}

	if len(cfg.DefaultPSK) == 0 {
		return Config{}, fmt.Errorf("MESHTASTIC_DEFAULT_PSK não pode ser vazio")
	}

	return cfg, nil
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
