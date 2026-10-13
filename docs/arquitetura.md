# Arquitetura do DevNux Mesh Brasil

## Objetivo

Este documento descreve a organização do código e o fluxo de uma publicação MQTT. O modelo de protocolo está em [Protocolo e criptografia](protocolo.md), a operação em [Configuração](configuracao.md), o processo de manutenção em [Guia de desenvolvimento](desenvolvimento.md) e os resultados físicos em [Equipamentos e ensaios](equipamentos.md) e [Validação](validacao.md).

## Organização

- cmd/broker: composição da aplicação, configuração, logger, servidor MQTT, TLS e encerramento.
- internal/broker: integração com eventos MQTT e processamento das publicações Meshtastic.
- internal/config: carregamento e validação da configuração.
- internal/crypto: nonce e transformação AES-CTR utilizados no processamento atual.

## Fluxo de uma publicação

1. O cliente estabelece conexão MQTT sobre TLS.
2. O broker recebe uma publicação.
3. O hook MQTT entrega o payload ao processador.
4. O payload é interpretado como ServiceEnvelope protobuf.
5. O envelope é validado.
6. O pacote é descriptografado usando a PSK configurada.
7. O resultado é interpretado como Data protobuf.
8. O fluxo registra os dados relevantes.
9. Erros de parsing, validação ou descriptografia interrompem o processamento daquela publicação.

## Limite importante

O hook atual tenta processar cada publicação como um envelope Meshtastic e usa uma PSK configurada. Publicações que não correspondam ao formato esperado ou não possam ser descriptografadas são ignoradas. Portanto, a implementação atual não deve ser descrita como encaminhador MQTT opaco e universal para quaisquer canais ou PSKs.

O processamento criptográfico utiliza AES-CTR, que não autentica nem protege a integridade do payload; não se deve afirmar que toda PSK incorreta será detectada. O transporte TLS protege a conexão, mas não fornece por si só autenticação MQTT de clientes nem autorização por tópico. A configuração atual não implementa esses controles.

## Princípio arquitetural

Evitar camadas artificiais. Separar responsabilidades quando isso melhora a compreensão ou permite testar uma regra isoladamente. Não adicionar interfaces genéricas, camada de repositório, service layer, injeção de dependência ou event bus sem necessidade demonstrada.

## Testes

Os testes unitários cobrem responsabilidades do hook, processador, configuração e criptografia. A cobertura descrita nos testes não substitui uma execução MQTT com dois clientes conectados nem um ensaio físico com dois ou mais rádios. Consulte [Validação](validacao.md) para o estado documentado e [Equipamentos e ensaios](equipamentos.md) para a matriz de hardware.

Qualquer mudança estrutural deve preservar o comportamento validado e ser justificada pelo requisito novo.