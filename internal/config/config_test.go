package config

import (
	"bytes"
	"testing"
)

func TestDefaultPSK(t *testing.T) {
	expected := []byte{
		0xd4, 0xf1, 0xbb, 0x3a,
		0x20, 0x29, 0x07, 0x59,
		0xf0, 0xbc, 0xff, 0xab,
		0xcf, 0x4e, 0x69, 0x01,
	}

	got := defaultPSK()
	if !bytes.Equal(got, expected) {
		t.Fatalf("PSK padrão incompatível: %x != %x", got, expected)
	}
}

func TestLoadUsesDefaultPSK(t *testing.T) {
	t.Setenv("MESHTASTIC_DEFAULT_PSK", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() retornou erro: %v", err)
	}

	if !bytes.Equal(cfg.DefaultPSK, defaultPSK()) {
		t.Fatalf("Load() não usou a PSK padrão: %x != %x", cfg.DefaultPSK, defaultPSK())
	}

	if cfg.ListenAddress != ":8883" {
		t.Fatalf("ListenAddress inesperado: %q", cfg.ListenAddress)
	}

	if cfg.CertificateFile != "certificate.pem" {
		t.Fatalf("CertificateFile inesperado: %q", cfg.CertificateFile)
	}

	if cfg.PrivateKeyFile != "private.key" {
		t.Fatalf("PrivateKeyFile inesperado: %q", cfg.PrivateKeyFile)
	}
}

func TestLoadUsesEnvironmentPSK(t *testing.T) {
	const customPSK = "0123456789abcdef"

	t.Setenv("MESHTASTIC_DEFAULT_PSK", customPSK)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() retornou erro: %v", err)
	}

	if !bytes.Equal(cfg.DefaultPSK, []byte(customPSK)) {
		t.Fatalf("PSK de ambiente inesperada: %x != %x", cfg.DefaultPSK, []byte(customPSK))
	}
}

func TestEnvOrDefault(t *testing.T) {
	t.Setenv("DEVNUX_TEST_ENV_OR_DEFAULT", "")
	if got := envOrDefault("DEVNUX_TEST_ENV_OR_DEFAULT", "fallback"); got != "fallback" {
		t.Fatalf("fallback inesperado: %q", got)
	}

	t.Setenv("DEVNUX_TEST_ENV_OR_DEFAULT", "configured")
	if got := envOrDefault("DEVNUX_TEST_ENV_OR_DEFAULT", "fallback"); got != "configured" {
		t.Fatalf("valor configurado inesperado: %q", got)
	}
}
