# DevNux Mesh Brasil

[![CI](https://github.com/marcelositr/devnux-mesh-brasil/actions/workflows/ci.yml/badge.svg)](https://github.com/marcelositr/devnux-mesh-brasil/actions/workflows/ci.yml)

Broker MQTT para integração de redes Meshtastic.

O **DevNux Mesh Brasil** é um projeto de infraestrutura de código aberto voltado ao recebimento e processamento de mensagens Meshtastic por MQTT. O projeto implementa o transporte MQTT sobre TLS, valida envelopes protobuf, descriptografa pacotes Meshtastic e registra os dados processados.

## Estado atual

O projeto está em desenvolvimento e validação técnica. O núcleo atual está coberto por testes automatizados e foi validado com `go test`, `go vet` e compilação do broker.

O escopo atual é deliberadamente pequeno: o objetivo é manter um broker funcional, compreensível e fácil de evoluir sem introduzir abstrações desnecessárias.

## O que o broker faz

- inicia um broker MQTT sobre TLS;
- aceita conexões e inscrições MQTT;
- recebe `ServiceEnvelope` codificados em protobuf;
- valida os campos necessários para o processamento;
- descriptografa pacotes Meshtastic com AES-CTR;
- interpreta o conteúdo descriptografado como `Data`;
- registra mensagens de texto e outros portnums reconhecidos pelo fluxo atual;
- encerra o processo de forma controlada por SIGINT ou SIGTERM.

A arquitetura e o fluxo interno estão descritos em [Arquitetura](docs/arquitetura.md).

## Configuração

Valores padrão:

| Configuração | Padrão |
|---|---|
| Endereço MQTT | `:8883` |
| Certificado TLS | `certificate.pem` |
| Chave privada TLS | `private.key` |
| PSK padrão | `d4f1bb3a20290759f0bcffabcf4e6901` |

Variáveis de ambiente:

- `MQTT_LISTEN_ADDRESS`
- `MQTT_CERTIFICATE_FILE`
- `MQTT_PRIVATE_KEY_FILE`
- `MESHTASTIC_DEFAULT_PSK`

A PSK fornecida por ambiente deve ser hexadecimal e resultar em 16 ou 32 bytes.

Os detalhes estão em [Configuração e execução](docs/configuracao.md).

## Docker

A imagem do broker é construída em múltiplas etapas e executa o binário como usuário sem privilégios.

```bash
docker build -t devnux-mesh-brasil .
```

Para executar com TLS, monte um diretório contendo `certificate.pem` e `private.key`:

```bash
docker run --rm \
  -p 8883:8883 \
  -v "$PWD/certs:/etc/devnux/certs:ro" \
  devnux-mesh-brasil
```

A imagem utiliza:

- `/etc/devnux/certs/certificate.pem` para o certificado;
- `/etc/devnux/certs/private.key` para a chave privada;
- `:8883` como endereço MQTT padrão;
- `MESHTASTIC_DEFAULT_PSK` para substituir a PSK padrão quando necessário.

O certificado e a chave privada não fazem parte da imagem nem do repositório.

## TLS

O listener utiliza TLS 1.2 como versão mínima e recebe certificado e chave privada de arquivos separados.

A configuração de TLS, os requisitos de execução e os cuidados operacionais estão documentados em [Configuração e execução](docs/configuracao.md).

## Protocolo e criptografia

O processamento utiliza os tipos protobuf do ecossistema Meshtastic e AES-CTR para a descriptografia dos pacotes.

As regras de validação, construção do nonce, uso da PSK e comportamento do processamento estão em [Protocolo e criptografia](docs/protocolo.md).

## Desenvolvimento

Para trabalhar no projeto:

```bash
go test ./...
go vet ./...
go build ./cmd/broker
```

As responsabilidades dos pacotes, fluxo de execução e critérios para alterações estão em [Guia de desenvolvimento](docs/desenvolvimento.md).

## Uso de IA

Ferramentas de IA podem ser utilizadas como apoio ao desenvolvimento, revisão, investigação e documentação. Elas não substituem a responsabilidade humana pelas decisões do projeto, pela validação do código ou pela verificação de resultados.

As regras adotadas pelo projeto estão em [Uso de IA](docs/uso-de-ia.md).

## Limites atuais

Ainda não fazem parte do núcleo atual:

- autenticação MQTT;
- ACL ou autorização por tópico;
- rate limiting;
- deduplicação;
- persistência;
- seleção de PSK por `channel_id`;
- filtragem avançada de portnums;
- observabilidade avançada;
- mecanismos de banimento ou moderação.

Esses limites são intencionais nesta etapa e não devem ser confundidos com funcionalidades já implementadas.

## Documentação

A documentação foi separada por responsabilidade para evitar duplicação:

- [Arquitetura](docs/arquitetura.md) — organização interna e fluxo do broker.
- [Protocolo e criptografia](docs/protocolo.md) — envelope, validação, nonce, AES-CTR e PSK.
- [Configuração e execução](docs/configuracao.md) — ambiente, TLS, execução e testes.
- [Guia de desenvolvimento](docs/desenvolvimento.md) — manutenção, mudanças e validação.
- [Uso de IA](docs/uso-de-ia.md) — regras para uso responsável de ferramentas de IA no desenvolvimento.

Cada documento trata apenas do seu próprio assunto e aponta para os demais quando necessário.
