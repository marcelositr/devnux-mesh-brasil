# Equipamentos e ensaios

Este documento registra os equipamentos utilizados ou planejados para validação do DevNux Mesh Brasil. A inclusão na lista não significa que o equipamento tenha sido testado: cada estado descreve apenas o que foi documentado.

## Matriz de equipamentos

| Equipamento | Papel | Estado | Evidência e limite |
|---|---|---|---|
| Heltec Wi-Fi LoRa 32 V4, ESP32-S3R2 | Rádio Meshtastic usado no primeiro ensaio físico com o broker | Testado no cenário documentado | Conectou por MQTT/TLS, publicou tráfego processado pelo broker e voltou a produzir tráfego após reconexão do Wi-Fi. Foi usado somente um rádio; isso não valida uma malha com vários nós. |
| M5Stack Cardputer Adv K132-ADV com módulo LoRa SX1262 e GPS ATGM336H | Segundo equipamento planejado para ensaios comparativos | Planejado, ainda não validado neste projeto | A validação deve ocorrer depois de o equipamento estar disponível e configurado. Não afirmar compatibilidade com o broker antes do ensaio. |
| Outros rádios Meshtastic | Possíveis equipamentos de comparação | Não documentado como testado neste projeto | Experiências pessoais ou testes fora deste repositório não comprovam validação deste broker. Registrar modelo, firmware e evidências quando usados. |

## Ambiente do ensaio físico já documentado

- Firmware Meshtastic no Heltec: 2.7.26.54e0d8d.
- CLI Meshtastic: 2.7.11.
- Transporte: MQTT sobre TLS em rede local.
- Canal observado: LongFast, com uplink e downlink habilitados.
- Canais privados secundários: uplink e downlink desabilitados durante o ensaio.
- Relatório de mapa MQTT: desabilitado.
- O certificado era autoassinado e o cliente não validava a identidade do servidor.

## Como registrar um novo ensaio

Para cada execução, anotar:

1. data aproximada e objetivo do teste;
2. modelo e revisão de hardware, quando disponíveis;
3. firmware do rádio, versão da CLI e versão do broker/commit;
4. canal e opções MQTT relevantes, sem publicar chaves privadas ou credenciais;
5. topologia: um rádio, dois rádios, gateway ou malha maior;
6. procedimento e evidência observável, como logs, saída de comandos ou resultado de CI;
7. resultado: passou, falhou, inconclusivo ou não executado;
8. limitações e questões ainda abertas.

## Regras de evidência

- Equipamento adquirido ou planejado não equivale a equipamento testado.
- Teste de compilação não equivale a teste de comunicação por rádio.
- Conexão TLS sem validação de certificado não comprova autenticação do servidor.
- Teste com um nó não comprova encaminhamento em uma malha com múltiplos nós.
- Não inferir causa de mensagens repetidas sem comparar identificadores e campos dos pacotes.
- Manter resultados anteriores mesmo quando um novo ensaio corrigir ou ampliar a conclusão.

Os resultados técnicos e os próximos casos de teste estão em [Validação](validacao.md).