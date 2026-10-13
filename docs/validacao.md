# Validação

Este documento registra resultados conhecidos de testes automatizados, Docker e hardware físico. Um resultado observado não deve ser generalizado para cenários ainda não testados. O inventário de equipamentos, versões e estado de cada ensaio está em [Equipamentos e ensaios](equipamentos.md).

## Validação automatizada

A cobertura presente nos testes do repositório inclui:

- preservação e encaminhamento de payload protobuf inválido sem interromper o fluxo MQTT;
- processamento de ServiceEnvelope Meshtastic válido;
- validação de campos obrigatórios e rejeição de envelopes/pacotes incompletos;
- descriptografia de um pacote válido e interpretação de diferentes portnums;
- validação do layout do nonce e dos valores de PSK padrão;
- round-trip AES-CTR e vetor de interoperabilidade com a PSK padrão;
- rejeição de chave AES com tamanho não suportado;
- teste determinístico com PSK incorreta para um vetor conhecido, sem atribuir ao AES-CTR uma capacidade de autenticação que ele não possui;
- fuzz target do parser `ServiceEnvelope`, executado por 30 segundos no workflow separado;
- rejeição de payload MQTT vazio;
- encaminhamento de `Data` com `portnum = 0` ou payload vazio;
- aceitação pelo hook de portnum desconhecido quando o payload é válido;
- preservação do tópico e dos campos de payload no hook `OnPublish`, incluindo falhas de inspeção, com reescrita de `hop_limit` quando zero-hop está ativo;
- encaminhamento de tópicos JSON de downlink sem interpretá-los como `ServiceEnvelope`;
- teste de integração MQTT/TLS que confirma entrega de `ServiceEnvelope`, `portnum = 0` e tópico JSON a um assinante.

**Limites da cobertura:** o teste de PSK incorreta cobre apenas um vetor fixo e garante que a chave alternativa não recupera aquele `Data` conhecido. AES-CTR não autentica o conteúdo; uma chave errada pode produzir bytes que falham no parsing ou, em alguns casos, dados que ainda podem ser interpretados. Não existe garantia criptográfica de detecção de chave incorreta. O teste de integração conecta dois clientes MQTT ao servidor TLS de teste e verifica a entrega de publicações protobuf e JSON sem que a inspeção local filtre os pacotes.

No commit `fd049bed560d25f81d3fdcd80a06d78df1202740`, a [CI](https://github.com/marcelositr/devnux-mesh-brasil/actions/runs/38445828267) passou com `gofmt`, `go test ./...`, `go vet ./...`, `go test -race ./...`, `govulncheck` (nenhuma vulnerabilidade alcançável encontrada), build Go, teste de exclusões do contexto Docker, construção e scan da imagem, geração e upload do SBOM CycloneDX e smoke test TLS. O [CodeQL](https://github.com/marcelositr/devnux-mesh-brasil/actions/runs/38445838203) também passou. O workflow [Fuzz](https://github.com/marcelositr/devnux-mesh-brasil/actions/runs/38445828477) passou após aproximadamente 869 mil execuções do parser em 30 segundos, sem falha. A imagem-base Go e a imagem Alpine de runtime estão fixadas por digest SHA-256; a CI passou com esses digests. O SBOM CycloneDX é publicado como artefato da execução, com retenção de 30 dias. A validação física multi-nó continua pendente.

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
- validação em hardware real da política zero-hop implementada no broker;
- prontidão para exposição pública.

Autenticação MQTT, ACL, rate limiting, deduplicação, persistência, moderação e observabilidade avançada continuam sendo possibilidades futuras, não funcionalidades já implementadas. Consulte [Evolução futura e hardening](evolucao.md).
