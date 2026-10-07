# DevNux Mesh Brasil

Infraestrutura MQTT para integração e comunicação entre redes Meshtastic no Brasil.

## Status

**Em desenvolvimento e testes.**

## Escopo atual

- Broker MQTT.
- Endpoint MQTT sobre TLS na porta 8883.
- TLS 1.2 como protocolo mínimo.
- Aceitação de conexões MQTT.
- Aceitação de inscrições MQTT.
- Recepção de `ServiceEnvelope` em protobuf.
- Validação dos campos essenciais do envelope.
- Descriptografia AES-CTR de pacotes Meshtastic.
- Decodificação do payload como `Data`.
- Registro de mensagens de texto e outros portnums.
- Encerramento controlado por SIGINT/SIGTERM.

## Configuração

Por padrão:

- endereço MQTT: `:8883`
- certificado: `certificate.pem`
- chave privada: `private.key`
- PSK padrão Meshtastic: `d4f1bb3a20290759f0bcffabcf4e6901` (16 bytes)

As configurações podem ser alteradas por variáveis de ambiente:

- `MQTT_LISTEN_ADDRESS`
- `MQTT_CERTIFICATE_FILE`
- `MQTT_PRIVATE_KEY_FILE`
- `MESHTASTIC_DEFAULT_PSK`

A PSK configurada por ambiente deve ser hexadecimal e possuir 16 ou 32 bytes após a decodificação.

Nunca coloque uma chave privada real no repositório.

## TLS

O broker utiliza certificado e chave privada em arquivos separados. O listener aceita somente TLS 1.2 ou superior.

Exemplo:

```bash
MQTT_CERTIFICATE_FILE=certificate.pem \
MQTT_PRIVATE_KEY_FILE=private.key \
./broker
```

## Testes locais

Dependências:

- Go instalado.
- Um certificado TLS para teste.
- Cliente MQTT para teste.

Comandos:

```bash
go test ./...
go vet ./...
go build ./cmd/broker
```

## Comportamento do broker

Conexões e inscrições MQTT são aceitas sem autenticação ou filtragem adicional nesta etapa.

Publicações são interpretadas como `ServiceEnvelope`. Payloads vazios, protobufs inválidos, envelopes inválidos e pacotes que não possam ser descriptografados ou interpretados são ignorados pelo fluxo de publicação.

O nome do tópico não possui filtragem específica nesta etapa. Ele é utilizado para registro das mensagens.

A descriptografia utiliza a PSK padrão configurada para o broker. Não há seleção de PSK por `channel_id`.

## Documentação técnica

Os detalhes do protocolo processado, validações, criptografia e decisões de implementação estão em:

`docs/reproducao-meshtastic-mqtt.md`

## Fora do escopo atual

Não fazem parte da implementação atual:

- rate limiting;
- deduplicação de pacotes;
- persistência;
- ACL ou autorização por tópico;
- filtragem de portnums;
- observabilidade avançada;
- mecanismos de banimento ou moderação.
