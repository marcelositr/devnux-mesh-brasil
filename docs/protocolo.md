# Protocolo e criptografia

## Objetivo

Registrar o comportamento de protocolo implementado pelo DevNux Mesh Brasil sem duplicar as informações de arquitetura e operação.

A estrutura interna está em [Arquitetura](arquitetura.md). Os valores configuráveis estão em [Configuração e execução](configuracao.md).

## Entrada MQTT

O payload de uma publicação MQTT é tratado como um `ServiceEnvelope` protobuf.

O nome do tópico não determina, nesta etapa, uma regra adicional de filtragem. O tópico pode ser registrado pelo fluxo MQTT, mas não é usado para selecionar a PSK.

## Validação do envelope

Antes da descriptografia, o processador exige:

- `channel_id` preenchido;
- `gateway_id` preenchido;
- `packet` presente;
- identificador do pacote válido;
- origem válida;
- payload criptografado presente;
- ausência de payload já decodificado no `MeshPacket`.

Falhas de validação fazem a publicação ser ignorada pelo fluxo de processamento.

## Criptografia

O processamento utiliza AES-CTR.

A chave aceita pelo núcleo criptográfico possui 16 ou 32 bytes, correspondendo aos tamanhos AES-128 e AES-256.

A PSK padrão do projeto é:

```
d4 f1 bb 3a 20 29 07 59 f0 bc ff ab cf 4e 69 01
```

A PSK de ambiente é interpretada como hexadecimal e deve produzir 16 ou 32 bytes.

## Nonce

O nonce possui 16 bytes.

Os oito primeiros bytes representam o identificador do pacote em little-endian. Os quatro bytes seguintes representam a origem em little-endian. Os quatro bytes finais permanecem zerados.

A mesma transformação AES-CTR é usada para criptografar e descriptografar dados.

## Após a descriptografia

O resultado é interpretado como `meshtastic.Data`.

O fluxo atual consegue registrar mensagens de texto e outros portnums conforme os dados presentes. Pacotes que não produzam os dados necessários para um registro textual não devem ser confundidos com falhas do broker.

## Compatibilidade e validação

O núcleo criptográfico possui teste de interoperabilidade baseado em vetor conhecido. A validação automatizada deve ser executada antes de alterações no protocolo ou na criptografia.

A validação com hardware Meshtastic físico e real ainda é uma etapa pendente. Qualquer decisão futura que altere filtragem, autorização ou controle de tráfego deve considerar o comportamento observado nessa etapa.

Os comandos de validação estão em [Guia de desenvolvimento](desenvolvimento.md).

## Limites

A implementação atual utiliza uma única PSK padrão para o broker. Não existe seleção automática de PSK por `channel_id`.

Também não há, nesta etapa, filtragem avançada de tópicos ou portnums.

Possíveis evoluções de hardening estão documentadas separadamente em [Evolução futura e hardening](evolucao.md) e não constituem requisitos da implementação atual.
