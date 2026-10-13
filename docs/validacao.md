# Validação

Este documento registra resultados conhecidos de testes automatizados, Docker e hardware físico. Um resultado observado não deve ser generalizado para cenários ainda não testados. O inventário de equipamentos, versões e estado de cada ensaio está em [Equipamentos e ensaios](equipamentos.md).

## Validação automatizada

A cobertura presente nos testes do repositório inclui:

- rejeição de payload protobuf inválido pelo hook;
- processamento de ServiceEnvelope Meshtastic válido;
- validação de campos obrigatórios e rejeição de envelopes/pacotes incompletos;
- descriptografia de um pacote válido e interpretação de diferentes portnums;
- validação do layout do nonce e dos valores de PSK padrão;
- round-trip AES-CTR e vetor de interoperabilidade com a PSK padrão;
- rejeição de chave AES com tamanho não suportado;
- rejeição de payload MQTT vazio;
- descarte pelo hook de `Data` com `portnum = 0` ou payload vazio;
- aceitação pelo hook de portnum desconhecido quando o payload é válido;
- preservação do payload original no retorno do hook `OnPublish` nos cenários cobertos.

**Limites da cobertura:** não existe um teste específico que comprove a rejeição de uma PSK incorreta. AES-CTR não autentica o conteúdo; uma chave errada pode produzir bytes que falham no parsing, mas não existe garantia criptográfica de detecção de chave incorreta. Há um teste de integração que conecta dois clientes MQTT ao servidor TLS de teste e verifica a entrega de uma publicação aceita. O teste também verifica que uma publicação com `portnum = 0` não chega ao assinante; a execução após essa alteração deve ser confirmada na CI. Na execução do commit `75b59f7`, a CI passou em `go test ./...`, `go vet ./...`, build Go, build Docker e teste de inicialização do container com validação TLS.

Esta lista descreve a cobertura encontrada no código, não uma execução realizada durante esta revisão. Reexecute `go test ./...` após alterações e consulte os resultados atuais da CI antes de afirmar que a suíte está passando.

## Validação Docker

O histórico de validação do projeto registra uma construção local com o comando **docker build -t devnux-mesh-brasil .**.

A imagem iniciou como usuário sem privilégios. Durante o ensaio inicial, a chave privada montada com permissão 0600 não pôde ser lida pelo usuário do container; o teste foi repetido ajustando o acesso de grupo sem remover o usuário não-root.

Também foi registrado um teste de publicação MQTT/TLS com um ServiceEnvelope válido, processado e descriptografado no container. O log de exemplo registrado na ocasião foi, antes da alteração que removeu o texto das mensagens dos logs:

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

A validação com múltiplos equipamentos fica **adiada até depois da conclusão do trabalho de software**, por disponibilidade limitada de rádios. Isso não impede concluir as correções, testes automatizados, documentação e verificações de build/CI que podem ser realizadas agora.

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
