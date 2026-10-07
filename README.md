# DevNux Mesh Brasil

Infraestrutura MQTT comunitária para integração e comunicação entre redes Meshtastic no Brasil.

## Reprodução do projeto de referência

Esta branch contém a reprodução inicial, em Go, do projeto meshtastic/mqtt.

O objetivo desta etapa é reproduzir o comportamento existente no projeto de referência antes de introduzir decisões específicas de arquitetura do DevNux.

## Status

**Em desenvolvimento e testes.**

## Escopo desta etapa

- Broker MQTT.
- Endpoint MQTT sobre TLS na porta 8883.
- Aceitação de conexões.
- Aceitação de inscrições MQTT.
- Recepção de ServiceEnvelope em protobuf.
- Validação do envelope.
- Descriptografia AES-CTR de pacotes de canal.
- Decodificação do payload Data.
- Registro de mensagens de texto e outros portnums.
- Encerramento controlado por SIGINT/SIGTERM.

## Certificados

O projeto de referência utiliza um certificado PFX com uma senha embutida no código. Isso não será reproduzido.

Nesta implementação, o certificado e a chave privada são arquivos separados e configuráveis por variáveis de ambiente:

- MQTT_CERTIFICATE_FILE
- MQTT_PRIVATE_KEY_FILE

Nunca coloque uma chave privada real no repositório.

## Configuração

Por padrão:

- MQTT TLS: :8883
- certificado: certificate.pem
- chave privada: private.key
- PSK padrão Meshtastic: `d4f1bb3a20290759f0bcffabcf4e6901` (16 bytes)

As configurações podem ser alteradas por variáveis de ambiente.

## Testes locais

Dependências:

- Go instalado.
- Um certificado TLS para teste.
- Cliente MQTT para teste.

Comandos:

    go mod tidy
    go test ./...
    go vet ./...
    go build ./cmd/broker

Para iniciar:

    MQTT_CERTIFICATE_FILE=certificate.pem \
    MQTT_PRIVATE_KEY_FILE=private.key \
    ./broker

O broker escuta somente em TLS, como o código atual do projeto de referência.

## Documentação

A análise de compatibilidade e os pontos que ainda precisam de validação estão em:

docs/reproducao-meshtastic-mqtt.md

## Observação

Esta branch ainda não representa a arquitetura final do DevNux Mesh Brasil. Ela existe para reproduzir e testar o comportamento do projeto de referência.
