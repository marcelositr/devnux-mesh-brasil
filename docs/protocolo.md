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

```text
d4 f1 bb 3a 20 29 07 59 f0 bc ff ab cf 4e 69 01
```

A PSK de ambiente é interpretada como hexadecimal e deve produzir 16 ou 32 bytes.

## Nonce

O nonce possui 16 bytes.

Os oito primeiros bytes representam o identificador do pacote em little-endian. Os quatro bytes seguintes representam a origem em little-endian. Os quatro bytes finais permanecem zerados.

A mesma transformação AES-CTR é usada para criptografar e descriptografar dados.

**Limite de segurança:** AES-CTR não fornece autenticação nem integridade do conteúdo. O teste com PSK incorreta usa um vetor conhecido para verificar um resultado específico; ele não transforma a cifra em um mecanismo de autenticação nem garante que todas as chaves incorretas sejam detectadas.

## Após a descriptografia

O resultado é interpretado como `meshtastic.Data`.

O fluxo registra mensagens de texto e outros portnums quando o `Data` contém payload. Portnums desconhecidos não são descartados apenas por serem desconhecidos, desde que o payload seja válido. `Data` com `portnum = 0` ou payload vazio é ignorado pelo hook e não é encaminhado aos assinantes MQTT. Falhas de parsing, validação ou descriptografia também são ignoradas.

## MQTT, hops e repetições

O hop limit do Meshtastic limita quantos encaminhamentos um pacote pode realizar na malha de rádios. Isso não determina quantas publicações MQTT o broker deve receber. Uma contagem de logs repetidos, isoladamente, não prova que sejam cópias do mesmo pacote nem que a causa seja o hop limit.

Para investigar repetições, é necessário comparar identificadores e campos do protocolo e observar o tráfego em uma malha com mais de um nó. O broker não implementa deduplicação nesta etapa; não adicionar essa regra sem evidência de um problema concreto e sem critérios compatíveis com o comportamento real da rede.

A configuração MQTT do Meshtastic também define quais canais fazem uplink para o servidor e quais recebem downlink. Consulte a [documentação oficial do módulo MQTT](https://meshtastic.org/docs/configuration/module/mqtt/) para a semântica das opções do rádio.

## Compatibilidade e validação

O núcleo criptográfico possui teste de interoperabilidade baseado em vetor conhecido. A validação automatizada deve ser executada antes de alterações no protocolo ou na criptografia.

Um teste físico inicial confirmou conexão MQTT/TLS de um rádio, recebimento e processamento de pacotes e chegada de mensagens após a reconexão do Wi-Fi. Ainda não foi validada uma malha física com múltiplos nós, nem foram concluídas as verificações de encaminhamento e downlink nessa topologia. Os resultados detalhados estão em [Validação](validacao.md).

Qualquer decisão futura que altere filtragem, autorização ou controle de tráfego deve considerar o comportamento observado em hardware real.

## Limites

A implementação atual utiliza uma única PSK padrão para o broker. Não existe seleção automática de PSK por `channel_id`.

Também não há, nesta etapa, filtragem avançada de tópicos ou portnums nem deduplicação no broker.

Possíveis evoluções de hardening estão documentadas separadamente em [Evolução futura e hardening](evolucao.md) e não constituem requisitos da implementação atual.
