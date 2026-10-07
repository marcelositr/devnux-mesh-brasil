package main

import (
	"crypto/tls"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"

	"github.com/marcelositr/devnux-mesh-brasil/internal/broker"
	"github.com/marcelositr/devnux-mesh-brasil/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("falha ao carregar configuração", "erro", err)
		os.Exit(1)
	}

	server := mqtt.New(nil)

	// O projeto de referência aceita conexões sem autenticação adicional.
	if err := server.AddHook(new(auth.AllowHook), nil); err != nil {
		logger.Error("falha ao configurar autorização MQTT", "erro", err)
		os.Exit(1)
	}

	handler := broker.NewHandler(logger, cfg)
	if err := server.AddHook(handler, nil); err != nil {
		logger.Error("falha ao configurar processamento de pacotes", "erro", err)
		os.Exit(1)
	}

	tlsConfig, err := loadTLSConfig(cfg)
	if err != nil {
		logger.Error("falha ao carregar certificado TLS", "erro", err)
		os.Exit(1)
	}

	listener := listeners.NewTCP(listeners.Config{
		ID:        "mqtt-tls",
		Address:   cfg.ListenAddress,
		TLSConfig: tlsConfig,
	})

	if err := server.AddListener(listener); err != nil {
		logger.Error("falha ao adicionar listener MQTT", "erro", err)
		os.Exit(1)
	}

	go func() {
		if err := server.Serve(); err != nil {
			logger.Error("servidor MQTT encerrado com erro", "erro", err)
		}
	}()

	logger.Info("servidor MQTT iniciado", "endereco", cfg.ListenAddress)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	logger.Info("sinal de encerramento recebido")
	if err := server.Close(); err != nil {
		logger.Error("falha ao encerrar servidor MQTT", "erro", err)
		os.Exit(1)
	}
	logger.Info("servidor MQTT encerrado")
}

func loadTLSConfig(cfg config.Config) (*tls.Config, error) {
	certificate, err := tls.LoadX509KeyPair(cfg.CertificateFile, cfg.PrivateKeyFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{certificate},
	}, nil
}
