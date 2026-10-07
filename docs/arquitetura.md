# Arquitetura do DevNux Mesh Brasil

## Objetivo

Este documento descreve como o código está organizado e como uma publicação MQTT percorre o broker.

Detalhes de protocolo e criptografia ficam em [Protocolo e criptografia](protocolo.md). Configuração e operação ficam em [Configuração e execução](configuracao.md). Regras de manutenção ficam em [Guia de desenvolvimento](desenvolvimento.md).

## Organização

```
cmd/broker/
    main.go

internal/
    broker/
        handler.go
        packet_processor.go
    config/
        config.go
    crypto/
        crypto.go
```

### `cmd/broker`

Responsável pela composição da aplicação:

- cria o logger;
- carrega a configuração;
- cria o broker MQTT;
- registra hooks;
- configura TLS;
- inicia o servidor;
- trata o encerramento do processo.

O pacote principal não contém regras do protocolo Meshtastic.

### `internal/broker`

Contém a integração entre MQTT e o processamento Meshtastic.

`handler.go` atua como adaptador dos eventos MQTT.

`packet_processor.go` concentra o processamento do payload Meshtastic: parsing, validação e descriptografia.

### `internal/config`

Carrega configurações do ambiente, aplica valores padrão e valida a PSK fornecida por variável de ambiente.

### `internal/crypto`

Contém a implementação criptográfica necessária ao processamento atual, incluindo a PSK padrão, construção do nonce e transformação AES-CTR.

## Fluxo de uma publicação

1. O cliente estabelece uma conexão MQTT sobre TLS.
2. O broker recebe uma publicação.
3. O hook MQTT entrega o payload ao processador.
4. O processador interpreta o payload como `ServiceEnvelope`.
5. O envelope é validado.
6. O pacote criptografado é descriptografado.
7. O resultado é interpretado como `Data`.
8. O fluxo registra os dados relevantes.
9. Erros de parsing, validação ou descriptografia interrompem o processamento daquela publicação sem derrubar o broker.

## Princípio arquitetural

O projeto evita camadas artificiais. Uma responsabilidade deve ser separada quando a separação melhora a compreensão ou permite testar uma regra isoladamente.

Não são necessários, no estado atual:

- interfaces para cada componente;
- camada de repositório;
- service layer genérica;
- container de injeção de dependência;
- event bus;
- abstrações sobre MQTT ou TLS sem necessidade concreta.

Qualquer mudança estrutural deve preservar o comportamento validado e ser justificada pelo novo requisito.

## Testes

Os testes acompanham as responsabilidades:

- `internal/broker/handler_test.go` testa a integração do hook MQTT;
- `internal/broker/packet_processor_test.go` testa o processamento Meshtastic;
- `internal/config/config_test.go` testa carregamento e validação da configuração;
- `internal/crypto/crypto_test.go` testa nonce, chaves e interoperabilidade criptográfica.

A validação completa está descrita em [Guia de desenvolvimento](desenvolvimento.md).
