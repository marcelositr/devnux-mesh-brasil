# Arquitetura conceitual

## Estado

Este documento descreve os limites conceituais atuais; não aprova funcionalidades futuras. A implementação concreta está descrita em [Arquitetura técnica](../docs/arquitetura.md).

## Componentes

- Composição da aplicação: configuração, logger, servidor MQTT, listener TLS e encerramento.
- Adaptador MQTT: recebe eventos de conexão, inscrição e publicação.
- Processador Meshtastic: interpreta e valida ServiceEnvelope e tenta descriptografar o pacote.
- Configuração: carrega endereço, caminhos de certificado e PSK.
- Criptografia: implementa transformação AES-CTR e construção de nonce necessárias ao processamento atual.

## Limite de responsabilidade

O hook atual tenta interpretar cada publicação como envelope Meshtastic antes de aceitá-la. Publicações incompatíveis com o formato ou com a PSK configurada podem ser ignoradas; o comportamento atual não deve ser descrito como encaminhamento opaco universal.

## Princípios estruturais

- Manter componentes pequenos e testáveis.
- Não introduzir camadas genéricas sem necessidade concreta.
- Não confundir transporte TLS com autenticação de cliente.
- Não adicionar deduplicação, seleção de PSK por canal ou políticas de tópicos sem requisito e plano de validação.
- Alterações de protocolo devem preservar testes de interoperabilidade e ser validadas com hardware quando necessário.

## Documentos relacionados

- [Arquitetura técnica](../docs/arquitetura.md)
- [Protocolo e criptografia](../docs/protocolo.md)
- [Validação](../docs/validacao.md)
- [Equipamentos e ensaios](../docs/equipamentos.md)