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
2. O broker recebe uma publicação e aplica o roteamento MQTT normal conforme os filtros dos assinantes.
3. O hook reconhece tópicos protobuf Meshtastic atuais (`/2/e/`) e legados (`/2/c/`) para inspeção opcional.
4. Quando o payload é um `ServiceEnvelope` que pode ser interpretado, o broker tenta registrar metadados do `Data`; o conteúdo textual não é gravado nos logs.
5. Falha de parsing, validação ou descriptografia afeta apenas a inspeção. O tópico e o payload originais continuam no fluxo MQTT.
6. Tópicos JSON, MapReport, status e outras famílias MQTT são encaminhados sem serem tratados como `ServiceEnvelope`.
7. Um gateway Meshtastic configurado para downlink assina os tópicos correspondentes e decide, segundo a configuração do rádio e do firmware, se encaminha os pacotes recebidos à malha.

## Limites importantes

O broker não cria pacotes LoRa nem transmite diretamente por USB/serial. O downlink físico depende de um nó gateway com o módulo MQTT habilitado e downlink ativado no canal apropriado. A PSK configurada no broker serve apenas à inspeção local de alguns pacotes; não deve ser usada como requisito para encaminhar tráfego cifrado, canais com chaves distintas, pacotes não decifráveis ou portnums ainda desconhecidos.

O broker preserva os tópicos e payloads recebidos, mas a configuração atual não implementa autenticação MQTT nem ACL por tópico. Por isso, deve permanecer em rede controlada até que os controles de acesso sejam projetados e validados.

O processamento criptográfico utiliza AES-CTR, que não autentica nem protege a integridade do payload; não se deve afirmar que toda PSK incorreta será detectada. O transporte TLS protege a conexão, mas não fornece por si só autenticação MQTT de clientes nem autorização por tópico. A configuração atual não implementa esses controles.

## Princípio arquitetural

Evitar camadas artificiais. Separar responsabilidades quando isso melhora a compreensão ou permite testar uma regra isoladamente. Não adicionar interfaces genéricas, camada de repositório, service layer, injeção de dependência ou event bus sem necessidade demonstrada.

## Testes

Os testes unitários cobrem responsabilidades do hook, processador, configuração e criptografia. A cobertura descrita nos testes não substitui uma execução MQTT com dois clientes conectados nem um ensaio físico com dois ou mais rádios. Consulte [Validação](validacao.md) para o estado documentado e [Equipamentos e ensaios](equipamentos.md) para a matriz de hardware.

Qualquer mudança estrutural deve preservar o comportamento validado e ser justificada pelo requisito novo.