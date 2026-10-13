# DevNux Mesh Brasil

[![CI](https://github.com/marcelositr/devnux-mesh-brasil/actions/workflows/ci.yml/badge.svg)](https://github.com/marcelositr/devnux-mesh-brasil/actions/workflows/ci.yml)

Broker MQTT para receber, validar e processar publicações Meshtastic por MQTT/TLS.

O **DevNux Mesh Brasil** é um projeto de infraestrutura de código aberto. O núcleo atual recebe payloads **ServiceEnvelope** em protobuf, valida os campos necessários, tenta descriptografar o pacote com uma PSK configurada e registra informações sobre os dados processados. O payload MQTT original é preservado no fluxo normal do broker quando a publicação é aceita pelo hook.

## Estado do projeto

O projeto está em desenvolvimento e validação técnica. A documentação distingue testes automatizados, testes em Docker e ensaios com hardware físico. Um teste com um rádio não comprova funcionamento em uma malha de múltiplos nós.

O teste físico documentado utilizou um **Heltec Wi-Fi LoRa 32 V4**. A validação com um segundo rádio, incluindo o **M5Stack Cardputer Adv**, continua pendente até que o ensaio seja realizado e seus resultados registrados. Consulte [Registro de equipamentos e ensaios](docs/equipamentos.md) e [Validação](docs/validacao.md).

## O que o broker faz atualmente

- inicia um broker MQTT sobre TLS;
- aceita conexões e inscrições MQTT;
- interpreta publicações como ServiceEnvelope protobuf;
- valida os campos necessários ao processamento;
- tenta descriptografar o pacote Meshtastic com AES-CTR usando a PSK padrão configurada;
- interpreta o resultado como Data protobuf;
- registra mensagens de texto e informações de outros portnums reconhecidos pelo fluxo atual;
- trata sinais de encerramento SIGINT e SIGTERM.

**Limite importante:** o processamento depende da PSK configurada. Publicações que não possam ser descriptografadas ou interpretadas pelo fluxo atual são ignoradas pelo hook. Portanto, não trate a implementação atual como um encaminhador transparente e universal para qualquer canal/PSK Meshtastic.

A estrutura interna está descrita em [Arquitetura](docs/arquitetura.md), e o comportamento do protocolo em [Protocolo e criptografia](docs/protocolo.md).

## Configuração

| Variável | Padrão |
|---|---|
| MQTT_LISTEN_ADDRESS | :8883 |
| MQTT_CERTIFICATE_FILE | certificate.pem |
| MQTT_PRIVATE_KEY_FILE | private.key |
| MESHTASTIC_DEFAULT_PSK | PSK padrão do projeto |

A PSK fornecida por ambiente deve ser hexadecimal e resultar em 16 ou 32 bytes. Veja [Configuração e execução](docs/configuracao.md).

## Docker

A imagem usa build multi-stage, fixa as imagens-base por digest SHA-256 e executa o binário como usuário sem privilégios. O certificado e a chave privada são fornecidos em tempo de execução, não embutidos na imagem.

Construção local: **docker build -t devnux-mesh-brasil .**

Imagem publicada no GitHub Container Registry: **ghcr.io/marcelositr/devnux-mesh-brasil:latest**

Exemplo limitado ao próprio host:

    docker run --rm -p 127.0.0.1:8883:8883 -v "$PWD/certs:/etc/devnux/certs:ro" ghcr.io/marcelositr/devnux-mesh-brasil:latest

Os caminhos padrão no container são /etc/devnux/certs/certificate.pem e /etc/devnux/certs/private.key.

**Segurança:** a configuração atual permite conexões MQTT sem autenticação e não implementa ACL por tópico. TLS protege o transporte, mas não substitui autenticação MQTT nem garante que o cliente valide a identidade do servidor. Não exponha o listener à Internet ou a redes não confiáveis.

## Desenvolvimento e testes

    go test ./...
    go vet ./...
    go test -race ./...
    go run "golang.org/x/vuln/cmd/govulncheck@$(tr -d '[:space:]' < .github/govulncheck-version)" ./...
    go build ./cmd/broker

Fuzzing direcionado do parser protobuf (30 segundos):

    go test ./internal/broker -fuzz=FuzzParseServiceEnvelope -fuzztime=30s -parallel=2

A CI verifica formatação, testes, `vet`, detector de condições de corrida, vulnerabilidades Go alcançáveis, build Go, exclusões de segredos no contexto Docker, construção e scan da imagem, geração de SBOM CycloneDX e inicialização TLS do container. O workflow [Fuzz](.github/workflows/fuzz.yml) executa fuzzing direcionado do parser `ServiceEnvelope`. As Actions de terceiros estão fixadas por SHA completo e as imagens-base por digest SHA-256. Confirme sempre os resultados na execução do commit candidato.

Consulte [Guia de desenvolvimento](docs/desenvolvimento.md) e [Validação](docs/validacao.md) para distinguir cobertura existente de execuções efetivamente confirmadas.

## Limites atuais

A implementação não inclui:

- autenticação MQTT de clientes;
- ACL/autorização por tópico;
- rate limiting configurado pelo projeto;
- deduplicação no broker;
- persistência própria;
- seleção de PSK por channel_id;
- filtragem avançada de tópicos ou portnums;
- moderação ou banimento de clientes;
- observabilidade avançada.

Esses itens não são compromissos automáticos de implementação. Devem ser reconsiderados com base em requisitos concretos, risco operacional e observações de hardware. Consulte [Evolução futura e hardening](docs/evolucao.md).

## Documentação

- [Arquitetura](docs/arquitetura.md): componentes e fluxo interno.
- [Protocolo e criptografia](docs/protocolo.md): envelope, validação, nonce, AES-CTR e PSK.
- [Configuração e execução](docs/configuracao.md): variáveis, TLS, Docker e cuidados operacionais.
- [Guia de desenvolvimento](docs/desenvolvimento.md): processo de alteração e validação.
- [Validação](docs/validacao.md): resultados e limitações conhecidos.
- [Equipamentos e ensaios](docs/equipamentos.md): matriz de hardware, ambiente e resultados físicos.
- [Evolução futura e hardening](docs/evolucao.md): possibilidades ainda não aprovadas como requisitos.
- [Uso de IA](docs/uso-de-ia.md): uso responsável de ferramentas de IA.
- [Governança](GOVERNANCE/README.md): princípios, decisões e regras de evolução.

## Identidade e governança

O nome canônico do repositório é **DevNux Mesh Brasil**. A descrição curta do projeto é: *broker MQTT para integração e processamento de publicações Meshtastic*. A governança orienta decisões; a documentação técnica descreve o comportamento implementado; o código é a implementação verificável desse comportamento.
