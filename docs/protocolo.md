# Protocolo e criptografia

## Objetivo

Registrar o comportamento de protocolo implementado pelo DevNux Mesh Brasil sem duplicar as informações de arquitetura e operação.

A estrutura interna está em [Arquitetura](arquitetura.md). Os valores configuráveis estão em [Configuração e execução](configuracao.md).

## Famílias de tráfego MQTT

O broker deve preservar o tópico e o payload originais. Nem todo tópico Meshtastic contém um `ServiceEnvelope`, e nem todo pacote é descriptografável com a PSK padrão.

- **Protobuf de malha:** tópicos atuais como `msh/REGION/2/e/CHANNELNAME/USERID`; firmwares anteriores a 2.3.0 podem usar `/2/c/`. O payload normalmente é um `ServiceEnvelope` protobuf que encapsula um `MeshPacket`.
- **JSON:** quando habilitado e suportado pelo firmware, tópicos como `msh/REGION/2/json/CHANNELNAME/USERID` carregam representações de um subconjunto de portnums. Algumas versões documentam o tópico especial `msh/REGION/2/json/mqtt/` para comandos de downlink, exigindo um canal chamado `mqtt` com downlink ativado.
- **MapReport e tópicos auxiliares:** relatórios de mapa e outras famílias de tópicos não devem ser interpretados como `ServiceEnvelope` apenas por usarem o mesmo broker MQTT.
- **Conteúdo de malha:** `MeshPacket` pode transportar mensagens de texto, posição, NodeInfo, telemetria, roteamento/ACK, waypoint, dados de administração e outros portnums. O broker não deve rejeitar portnums desconhecidos apenas por não os interpretar.

Os nomes de tópicos dependem do root topic, da região, do canal, da versão do firmware e das opções MQTT. A inspeção local do broker é apenas observabilidade; falhas de parsing, validação ou descriptografia não bloqueiam a publicação aos assinantes.

## Inspeção opcional de ServiceEnvelope

Para tópicos protobuf reconhecidos, o processador tenta interpretar o payload como `ServiceEnvelope`. O nome do tópico não seleciona a PSK, e a PSK local não representa necessariamente os canais privados ou outros modos de criptografia. A falha de inspeção não altera nem descarta o payload MQTT.

## Validação do envelope

Antes da descriptografia, o processador exige:

- `channel_id` preenchido;
- `gateway_id` preenchido;
- `packet` presente;
- identificador do pacote válido;
- origem válida;
- payload criptografado presente;
- ausência de payload já decodificado no `MeshPacket`.

Falhas de validação impedem apenas a inspeção do envelope; a publicação MQTT original continua sendo encaminhada.

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

Quando a inspeção consegue interpretar o `Data`, o fluxo registra metadados de mensagens de texto e outros portnums sem registrar o conteúdo textual. Portnums desconhecidos, `portnum = 0`, payload vazio e falhas de parsing, validação ou descriptografia não são motivo para bloquear a publicação MQTT original.

## MQTT, hops e repetições

O hop limit do Meshtastic limita quantos encaminhamentos um pacote pode realizar na malha de rádios. Isso não determina quantas publicações MQTT o broker deve receber. Uma contagem de logs repetidos, isoladamente, não prova que sejam cópias do mesmo pacote nem que a causa seja o hop limit.

Para investigar repetições, é necessário comparar identificadores e campos do protocolo e observar o tráfego em uma malha com mais de um nó. O broker não implementa deduplicação nesta etapa; não adicionar essa regra sem evidência de um problema concreto e sem critérios compatíveis com o comportamento real da rede. Encaminhar mensagens duplicadas é possível quando vários gateways publicam o mesmo pacote; consumidores que precisem deduplicar devem comparar o ID do pacote e outros campos de protocolo.

A configuração MQTT do Meshtastic também define quais canais fazem uplink para o servidor e quais recebem downlink. Consulte a [documentação oficial do módulo MQTT](https://meshtastic.org/docs/configuration/module/mqtt/) para a semântica das opções do rádio.

## Compatibilidade e validação

O núcleo criptográfico possui teste de interoperabilidade baseado em vetor conhecido. A validação automatizada deve ser executada antes de alterações no protocolo ou na criptografia.

Um teste físico inicial confirmou conexão MQTT/TLS de um rádio, recebimento e processamento de pacotes e chegada de mensagens após a reconexão do Wi-Fi. Ainda não foi validada uma malha física com múltiplos nós, nem foram concluídas as verificações de encaminhamento e downlink nessa topologia. Os resultados detalhados estão em [Validação](validacao.md).

Qualquer decisão futura que altere filtragem, autorização ou controle de tráfego deve considerar o comportamento observado em hardware real.

## Limites

A implementação atual utiliza uma única PSK padrão para o broker. Não existe seleção automática de PSK por `channel_id`.

Também não há, nesta etapa, filtragem avançada de tópicos ou portnums nem deduplicação no broker.

Possíveis evoluções de hardening estão documentadas separadamente em [Evolução futura e hardening](evolucao.md) e não constituem requisitos da implementação atual.
