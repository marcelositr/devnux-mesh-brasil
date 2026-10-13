# Protocolo e criptografia

## Objetivo

Registrar o comportamento de protocolo implementado pelo DevNux Mesh Brasil sem duplicar as informações de arquitetura e operação.

A estrutura interna está em [Arquitetura](arquitetura.md). Os valores configuráveis estão em [Configuração e execução](configuracao.md).

## Famílias de tráfego MQTT

O broker preserva o tópico e os campos do payload, exceto pela reescrita de `MeshPacket.hop_limit` para zero em envelopes protobuf reconhecidos quando `MQTT_ZERO_HOP=true` (padrão). Nem todo tópico Meshtastic contém um `ServiceEnvelope`, e nem todo pacote é descriptografável com a PSK padrão.

- **Protobuf de malha:** tópicos atuais como `msh/REGION/2/e/CHANNELNAME/USERID`; firmwares anteriores a 2.3.0 podem usar `/2/c/`. O payload normalmente é um `ServiceEnvelope` protobuf que encapsula um `MeshPacket`.
- **JSON:** quando habilitado e suportado pelo firmware, tópicos como `msh/REGION/2/json/CHANNELNAME/USERID` carregam representações de um subconjunto de portnums. Algumas versões documentam o tópico especial `msh/REGION/2/json/mqtt/` para comandos de downlink, exigindo um canal chamado `mqtt` com downlink ativado.
- **MapReport e tópicos auxiliares:** relatórios de mapa e outras famílias de tópicos não devem ser interpretados como `ServiceEnvelope` apenas por usarem o mesmo broker MQTT.
- **Conteúdo de malha:** `MeshPacket` pode transportar mensagens de texto, posição, NodeInfo, telemetria, roteamento/ACK, waypoint, dados de administração e outros portnums. O broker não deve rejeitar portnums desconhecidos apenas por não os interpretar.

Os nomes de tópicos dependem do root topic, da região, do canal, da versão do firmware e das opções MQTT. A inspeção local do broker é apenas observabilidade; falhas de parsing, validação ou descriptografia não bloqueiam a publicação aos assinantes. A reescrita zero-hop aplica-se somente quando o payload pode ser decodificado como `ServiceEnvelope` em um tópico protobuf reconhecido.

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

## Inventário de PortNums e payloads

O campo `Data.portnum` identifica o tipo de aplicação; o broker não deve exigir que a versão local das protobufs reconheça cada valor para encaminhar o pacote. O catálogo oficial evolui e inclui aplicações de terceiros e privadas. A tabela resume os valores registrados no [`portnums.proto` oficial](https://github.com/meshtastic/protobufs/blob/master/meshtastic/portnums.proto); ela é referência de protocolo, não uma lista de filtros do broker.

| PortNum | Nome | Categoria |
|---:|---|---|
| 0 | `UNKNOWN_APP` | Payload opaco legado |
| 1 | `TEXT_MESSAGE_APP` | Texto |
| 2 | `REMOTE_HARDWARE_APP` | Controle de hardware remoto |
| 3 | `POSITION_APP` | Posição |
| 4 | `NODEINFO_APP` | Identidade e informações do nó |
| 5 | `ROUTING_APP` | Roteamento, ACK/NACK e descoberta de rota |
| 6 | `ADMIN_APP` | Administração e configuração |
| 7 | `TEXT_MESSAGE_COMPRESSED_APP` | Texto comprimido |
| 8 | `WAYPOINT_APP` | Waypoints |
| 9 | `AUDIO_APP` | Quadros de áudio |
| 10 | `DETECTION_SENSOR_APP` | Sensor de detecção |
| 11 | `ALERT_APP` | Alertas |
| 12 | `KEY_VERIFICATION_APP` | Verificação de chaves |
| 13 | `REMOTE_SHELL_APP` | Shell remoto |
| 32 | `REPLY_APP` | Resposta/ping |
| 33 | `IP_TUNNEL_APP` | Túnel IP |
| 34 | `PAXCOUNTER_APP` | Contagem de pessoas |
| 35 | `STORE_FORWARD_PLUSPLUS_APP` | Store-and-forward++ |
| 36 | `NODE_STATUS_APP` | Estado do nó |
| 37 | `MESH_BEACON_APP` | Beacon de malha |
| 38 | `PAGING_APP` | Paging com confirmação humana |
| 64 | `SERIAL_APP` | Transporte serial |
| 65 | `STORE_FORWARD_APP` | Store-and-forward |
| 66 | `RANGE_TEST_APP` | Teste de alcance |
| 67 | `TELEMETRY_APP` | Telemetria |
| 68 | `ZPS_APP` | Estimativa de posição |
| 69 | `SIMULATOR_APP` | Simulador |
| 70 | `TRACEROUTE_APP` | Traceroute |
| 71 | `NEIGHBORINFO_APP` | Informações de vizinhos |
| 72 | `ATAK_PLUGIN` | Integração ATAK |
| 73 | `MAP_REPORT_APP` | Relatório de mapa |
| 74 | `POWERSTRESS_APP` | Teste de consumo de energia |
| 75 | `LORAWAN_BRIDGE` | Ponte LoRaWAN |
| 76 | `RETICULUM_TUNNEL_APP` | Túnel Reticulum |
| 77 | `CAYENNE_APP` | Cayenne Low Power Payload |
| 78 | `ATAK_PLUGIN_V2` | Integração ATAK v2 |
| 79 | `LORA_OTA_APP` | Atualização de firmware via LoRa |
| 112 | `GROUPALARM_APP` | Integração GroupAlarm |
| 256 | `PRIVATE_APP` | Identificador genérico para aplicações privadas |
| 257 | `ATAK_FORWARDER` | Integração ATAK Forwarder registrada |
| 256–511 | Faixa para aplicações privadas | O catálogo oficial também pode atribuir valores específicos dentro da faixa |

Além do `Data`, o `MeshPacket` carrega campos de encaminhamento e controle como origem, destino, ID do pacote, canal, limite inicial/atual de saltos, solicitação de ACK e próximo salto. O `ServiceEnvelope` acrescenta `channel_id` e `gateway_id`. O broker deve preservar o protobuf original em vez de reconstruir um pacote com campos que não compreende.

O MQTT também pode transportar `MapReport` e tópicos auxiliares fora do formato comum `ServiceEnvelope`. Esses payloads devem ser preservados, mesmo que o broker não os decodifique.

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

### Downlink físico e política zero-hop

O firmware Meshtastic assina os tópicos protobuf dos canais com `downlink_enabled` e encaminha os envelopes compatíveis para a malha. A recepção depende de `channel_id` corresponder ao canal configurado no rádio, do envelope ser válido e dos filtros de downlink do firmware; uma publicação MQTT aceita pelo broker não garante, por si só, transmissão LoRa.

O servidor MQTT público do Meshtastic aplica uma política de **zero-hop** para que o tráfego vindo da Internet alcance os nós diretamente conectados, sem se propagar indefinidamente pela malha local. O broker deste projeto agora aplica o mesmo princípio por padrão: `MQTT_ZERO_HOP=true` reescreve `MeshPacket.hop_limit` para `0` em envelopes protobuf reconhecidos e preserva `hop_start`. Defina `MQTT_ZERO_HOP=false` somente em uma malha privada e isolada onde a propagação multi-hop de tráfego MQTT seja intencional. A chave padrão do LongFast é pública; a documentação oficial recomenda canais privados com PSKs próprias para brokers privados. A política zero-hop reduz a propagação, mas não substitui autenticação MQTT, ACL nem validação física.

Para o firmware usado neste projeto, o caminho de downlink a validar primeiro é o protobuf `ServiceEnvelope`; não presumir que o tópico JSON de downlink documentado para versões compatíveis seja aceito por todo firmware. Referências: [MQTT Meshtastic](https://meshtastic.org/docs/software/integrations/mqtt/) e [configuração do módulo MQTT](https://meshtastic.org/docs/configuration/module/mqtt/).

## Compatibilidade e validação

O núcleo criptográfico possui teste de interoperabilidade baseado em vetor conhecido. A validação automatizada deve ser executada antes de alterações no protocolo ou na criptografia.

Um teste físico inicial confirmou conexão MQTT/TLS de um rádio, recebimento e processamento de pacotes e chegada de mensagens após a reconexão do Wi-Fi. Ainda não foi validada uma malha física com múltiplos nós, nem foram concluídas as verificações de encaminhamento e downlink nessa topologia. Os resultados detalhados estão em [Validação](validacao.md).

Qualquer decisão futura que altere filtragem, autorização ou controle de tráfego deve considerar o comportamento observado em hardware real.

## Limites

A implementação atual utiliza uma única PSK padrão para o broker. Não existe seleção automática de PSK por `channel_id`.

Também não há, nesta etapa, filtragem avançada de tópicos ou portnums nem deduplicação no broker.

Possíveis evoluções de hardening estão documentadas separadamente em [Evolução futura e hardening](evolucao.md) e não constituem requisitos da implementação atual.
