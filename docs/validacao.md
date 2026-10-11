# Validação

Este documento registra os resultados observados durante a validação do broker, distinguindo testes automatizados, testes em Docker e testes com hardware físico. Um resultado observado não deve ser generalizado para cenários que ainda não foram testados.

## Validação local automatizada

A suíte automatizada cobre:

- listener MQTT sobre TLS;
- publicação de payload protobuf inválido, rejeitado antes da distribuição;
- publicação de `ServiceEnvelope` Meshtastic válido;
- descriptografia com a PSK padrão;
- processamento de mensagens de texto;
- distribuição do payload original aos assinantes MQTT;
- validação de `ServiceEnvelope` sem `Packet`;
- validação de `MeshPacket` sem payload criptografado;
- rejeição de pacote com PSK incorreta;
- processamento de diferentes `portnum`;
- validação de nonce com valores de borda para `from` e `packetID`;
- rejeição de payload MQTT vazio.

Os comandos da bateria completa estão em [Guia de desenvolvimento](desenvolvimento.md). Os resultados dos testes automatizados devem ser confirmados novamente após alterações de código; esta lista descreve a cobertura existente, não substitui uma execução atual.

## Validação Docker

A imagem foi construída localmente com:

```bash
docker build -t devnux-mesh-brasil .
```

A imagem iniciou corretamente como usuário sem privilégios. Durante o primeiro teste, a chave privada montada com permissão `0600` não pôde ser lida pelo usuário do container. A execução foi então repetida com acesso de grupo à chave, preservando o usuário não-root da imagem.

O broker iniciou no container com:

```text
mochi mqtt server started
```

Em seguida, foi realizada uma publicação MQTT sobre TLS através da porta exposta pelo container. Um `ServiceEnvelope` Meshtastic válido foi processado e descriptografado corretamente.

Resultado observado no container:

```text
mensagem de texto recebida
topic=docker/teste
message=nonce test from=0x80000000 packetID=0xFFFFFFFF
```

Esse teste confirmou o fluxo principal no ambiente Docker:

```text
cliente MQTT
    -> TLS
    -> broker no container
    -> ServiceEnvelope
    -> descriptografia Meshtastic
    -> processamento
    -> publicação MQTT
```

A permissão da chave privada usada em testes deve continuar restrita. O certificado e a chave não devem ser adicionados ao repositório.

## Validação inicial com hardware físico

### Ambiente observado

- Rádio: Heltec Wi-Fi LoRa 32 V4 (ESP32-S3R2).
- Firmware Meshtastic: `2.7.26.54e0d8d`.
- CLI utilizada para inspeção: `2.7.11`.
- Transporte: MQTT sobre TLS para o broker em rede local.
- Canal principal observado: `LongFast`, com uplink e downlink habilitados.
- Os canais privados secundários estavam com uplink e downlink desabilitados.
- O relatório de mapa MQTT estava desabilitado.
- O rádio estava configurado para usar criptografia Meshtastic no MQTT.

### Resultados positivos

1. O broker escutou na porta TLS configurada e `openssl s_client` negociou TLS 1.3.
2. O rádio conectou ao broker e publicou tráfego Meshtastic.
3. O broker recebeu e processou pacotes, incluindo um pacote com `portnum=3` e mensagens de texto.
4. Após desativar e reativar o Wi-Fi do rádio, houve registros de reset/timeout durante a transição e, depois, novos pacotes e mensagens foram recebidos. Isso é evidência de recuperação da comunicação, mas ainda não é um teste prolongado de estabilidade.
5. A mensagem `Hi` apareceu em três registros separados por aproximadamente sete segundos.

### Limites e observações

- Foi testado somente um rádio físico. Não havia outros nós Meshtastic disponíveis para verificar encaminhamento entre rádios.
- A ocorrência de três registros de `Hi` foi observada, mas a causa não foi determinada. O limite de hops representa a possibilidade de encaminhamento na malha; não significa, por si só, que o broker MQTT deva receber uma cópia por hop. Os registros disponíveis não permitem concluir que as ocorrências foram duplicatas de um mesmo pacote.
- O teste não valida uma malha de múltiplos nós, a propagação de mensagens entre nós, a entrega de tráfego downlink a outro rádio ou o comportamento em uma rede com tráfego concorrente.
- O certificado usado no teste era autoassinado. O cliente do rádio estava configurado para não validar o certificado do servidor; portanto, o teste confirma a negociação TLS e a criptografia do transporte, mas não a autenticação da identidade do servidor.
- O broker não implementa autenticação MQTT nem ACL. Os testes devem permanecer em rede local controlada, sem exposição pública.

### Próximos testes com hardware

- [x] Conectar o rádio físico ao broker por MQTT/TLS.
- [x] Enviar mensagens e confirmar que o broker as processa.
- [x] Desativar e reativar o Wi-Fi e observar a chegada de tráfego após a reconexão.
- [ ] Repetir o teste de reconexão e observar a estabilidade por um período maior.
- [ ] Testar com pelo menos um segundo nó físico para observar encaminhamento real da malha.
- [ ] Comparar pacotes usando identificadores e campos do protocolo antes de classificar registros repetidos como duplicação.
- [ ] Verificar tráfego uplink/downlink com dois ou mais nós e documentar o resultado.
- [ ] Consolidar os resultados e revisar esta lista antes de considerar a fase de hardware concluída.

## Limites gerais da validação

Os testes automatizados e Docker cobrem o fluxo principal em ambiente controlado. O teste físico inicial confirma comunicação real de um rádio com o broker, mas **não confirma ainda o comportamento do broker em uma malha Meshtastic real com múltiplos nós**.

Autenticação MQTT, ACL, rate limiting, deduplicação no broker, persistência, moderação e observabilidade avançada continuam fora do escopo atual. Essas possibilidades estão descritas em [Evolução futura e hardening](evolucao.md) e devem ser reavaliadas somente após observar o comportamento real da malha.
