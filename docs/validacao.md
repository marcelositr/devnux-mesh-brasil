# Validação

Este documento registra resultados conhecidos de testes automatizados, Docker e hardware físico. Um resultado observado não deve ser generalizado para cenários ainda não testados. O inventário de equipamentos, versões e estado de cada ensaio está em [Equipamentos e ensaios](equipamentos.md).

## Validação automatizada

A documentação existente registra testes para:

- publicação de payload protobuf inválido;
- processamento de ServiceEnvelope Meshtastic válido;
- descriptografia com a PSK padrão e rejeição de PSK incorreta;
- processamento de mensagens de texto e diferentes portnums;
- validação de envelope sem Packet;
- validação de pacote sem payload criptografado;
- validação de nonce com valores de borda para from e packetID;
- rejeição de payload MQTT vazio;
- distribuição do payload original aos assinantes MQTT no cenário coberto pelos testes.

A lista descreve a cobertura documentada no repositório, não uma execução realizada durante esta revisão. Reexecute go test ./... após alterações e consulte os resultados atuais da CI antes de afirmar que a suíte está passando.

## Validação Docker

O histórico de validação do projeto registra uma construção local com o comando **docker build -t devnux-mesh-brasil .**.

A imagem iniciou como usuário sem privilégios. Durante o ensaio inicial, a chave privada montada com permissão 0600 não pôde ser lida pelo usuário do container; o teste foi repetido ajustando o acesso de grupo sem remover o usuário não-root.

Também foi registrado um teste de publicação MQTT/TLS com um ServiceEnvelope válido, processado e descriptografado no container. O log de exemplo registrado na ocasião foi:

    mensagem de texto recebida
    topic=docker/teste
    message=nonce test from=0x80000000 packetID=0xFFFFFFFF

Esse registro histórico não substitui uma nova execução depois de alterações no código, na imagem ou nas dependências.

## Validação inicial com hardware físico

### Ambiente documentado

- Rádio: Heltec Wi-Fi LoRa 32 V4, baseado em ESP32-S3R2.
- Firmware Meshtastic registrado no ensaio: 2.7.26.54e0d8d.
- CLI Meshtastic registrada: 2.7.11.
- Transporte: MQTT sobre TLS para broker em rede local.
- Canal observado: LongFast, com uplink e downlink habilitados.
- Canais privados secundários: uplink e downlink desabilitados durante o ensaio.
- Relatório de mapa MQTT: desabilitado.
- Criptografia Meshtastic para MQTT: habilitada.

### Resultados observados

1. Foi registrada negociação TLS 1.3 com openssl s_client.
2. O rádio conectou ao broker e publicou tráfego Meshtastic.
3. O broker recebeu e processou pacotes, incluindo um pacote com portnum=3 e mensagens de texto.
4. Após desativar e reativar o Wi-Fi do rádio, houve registros de reset/timeout durante a transição e depois novos pacotes e mensagens. Isso demonstra recuperação observada da comunicação naquele ensaio, não estabilidade prolongada.
5. A mensagem Hi apareceu em três registros separados por aproximadamente sete segundos. A causa não foi determinada.

### Limites do ensaio físico

- Apenas um rádio físico foi usado. Não foi validada uma malha com múltiplos nós.
- A ocorrência de três registros da mensagem não prova, por si só, duplicação do mesmo pacote. É necessário comparar identificadores e campos do protocolo.
- Não foi comprovada propagação entre rádios nem entrega de downlink a um segundo rádio.
- O certificado era autoassinado e o cliente não validava a identidade do servidor. O teste confirma a negociação TLS e o transporte cifrado, não a autenticação da identidade do servidor.
- O broker não implementa autenticação MQTT nem ACL. Os ensaios devem permanecer em rede local controlada.

## Próximos testes físicos

- [x] Conectar um rádio físico ao broker por MQTT/TLS.
- [x] Confirmar recebimento e processamento de mensagens.
- [x] Observar a chegada de tráfego após reconexão do Wi-Fi.
- [ ] Repetir reconexões e observar estabilidade por período maior.
- [ ] Testar com um segundo nó físico, incluindo o M5Stack Cardputer Adv quando disponível.
- [ ] Verificar uplink/downlink com dois ou mais nós.
- [ ] Comparar IDs e campos antes de classificar registros repetidos como duplicatas.
- [ ] Testar comportamento com PSKs/canais distintos e registrar quais publicações são aceitas ou ignoradas.
- [ ] Revisar este registro com evidências do ensaio antes de considerar concluída a fase de hardware.

## O que esta validação não comprova

A documentação disponível registra testes automatizados, um teste Docker e um primeiro teste físico com um rádio. Não comprova ainda:

- funcionamento em uma malha Meshtastic real com múltiplos nós;
- estabilidade prolongada;
- compatibilidade com todos os tipos de payload ou portnums;
- suporte transparente a canais com PSKs diferentes;
- autenticação de clientes MQTT ou autorização por tópico;
- prontidão para exposição pública.

Autenticação MQTT, ACL, rate limiting, deduplicação, persistência, moderação e observabilidade avançada continuam sendo possibilidades futuras, não funcionalidades já implementadas. Consulte [Evolução futura e hardening](evolucao.md).
