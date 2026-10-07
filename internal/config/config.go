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
		DefaultPSK:      defaultPSK(),
	}

	if value := os.Getenv("MESHTASTIC_DEFAULT_PSK"); value != "" {
		cfg.DefaultPSK = []byte(value)
	}

	if len(cfg.DefaultPSK) == 0 {
		return Config{}, fmt.Errorf("MESHTASTIC_DEFAULT_PSK não pode ser vazio")
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

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
