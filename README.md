# DevNux Mesh Brasil

Infraestrutura MQTT comunitária para integração e comunicação entre redes Meshtastic no Brasil.

## Reprodução do projeto de referência

Esta branch contém uma implementação em Go inspirada no comportamento do projeto meshtastic/mqtt. Diferenças deliberadas e limitações de validação estão registradas em `docs/reproducao-meshtastic-mqtt.md`.

## Status

**Em desenvolvimento e testes.** A validação de interoperabilidade com rádios físicos ainda é necessária.

## Escopo desta etapa

- Broker MQTT com endpoint TLS na porta 8883.
- TLS 1.2 como versão mínima.
- Aceitação de conexões e inscrições sem autenticação MQTT adicional.
- Inspeção opcional de `ServiceEnvelope` protobuf nos tópicos Meshtastic reconhecidos.
- Descriptografia AES-CTR para observabilidade local, sem bloquear o encaminhamento quando a inspeção falhar.
- Registro de metadados de mensagens sem registrar o texto recebido.
- Política zero-hop opcional.
- Encerramento controlado por SIGINT/SIGTERM.

## Certificados

O projeto de referência utiliza um certificado PFX com uma senha embutida no código. Isso não será reproduzido.

Nesta implementação, o certificado e a chave privada são arquivos separados e configuráveis por variáveis de ambiente:

- `MQTT_CERTIFICATE_FILE`
- `MQTT_PRIVATE_KEY_FILE`

Nunca coloque uma chave privada real no repositório.

## Configuração

Padrões:

- MQTT TLS: `:8883`
- certificado: `certificate.pem`
- chave privada: `private.key`
- PSK padrão Meshtastic: `d4f1bb3a20290759f0bcffabcf4e6901` (16 bytes)

Variáveis:

- `MQTT_LISTEN_ADDRESS`: endereço de escuta.
- `MQTT_CERTIFICATE_FILE`: caminho do certificado TLS.
- `MQTT_PRIVATE_KEY_FILE`: caminho da chave privada TLS.
- `MESHTASTIC_DEFAULT_PSK`: PSK em hexadecimal (compatibilidade anterior) ou Base64 padrão; após decodificar, deve ter 16 ou 32 bytes.
- `MESHTASTIC_ZERO_HOP`: `true` para definir `hop_limit=0) em envelopes protobuf válidos antes do encaminhamento; padrão `false`.

A PSK do canal Meshtastic é distinta das credenciais e do certificado usados por MQTT/TLS. Não grave uma PSK de implantação no repositório nem em logs.

## Testes locais

Dependências:

- Go instalado.
- Certificado e chave privada TLS para execução local.
- Cliente MQTT para teste.

Comandos:

    go test ./...
    go vet ./...
    go build -o /tmp/devnux-mesh-broker ./cmd/broker

Para iniciar:

    MQTT_CERTIFICATE_FILE=certificate.pem \
    MQTT_PRIVATE_KEY_FILE=private.key \
    /tmp/devnux-mesh-broker

O broker escuta somente em TLS.

## Documentação

A análise de compatibilidade, diferenças intencionais e pontos que ainda precisam de validação está em:

`docs/reproducao-meshtastic-mqtt.md`

Esta branch não representa a arquitetura final do DevNux Mesh Brasil. Ela existe para reproduzir, testar e validar o comportamento MQTT e a interoperabilidade com Meshtastic.
