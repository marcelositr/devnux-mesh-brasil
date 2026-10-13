# Reprodução do Meshtastic MQTT

## Objetivo

Esta branch implementa em Go o fluxo MQTT inspirado no projeto de referência `meshtastic/mqtt`. A inspeção local de pacotes é observacional: falhas de parsing ou descriptografia não devem interromper o encaminhamento MQTT.

## Comportamentos implementados

1. Broker MQTT com endpoint TLS na porta 8883.
2. TLS 1.2 como protocolo mínimo.
3. Conexões aceitas sem autenticação MQTT adicional.
4. Inscrições MQTT aceitas sem filtragem adicional.
5. Encaminhamento de publicações MQTT, inclusive JSON e payloads que não possam ser inspecionados localmente.
6. Inspeção de `ServiceEnvelope` protobuf nos tópicos Meshtastic reconhecidos (`msh/2/e/...` e `msh/2/c/...`).
7. Validação dos campos essenciais do envelope antes da inspeção criptográfica.
8. Descriptografia AES-CTR com nonce compatível com o formato Meshtastic implementado no pacote `internal/crypto`.
9. Decodificação do resultado como `Data` e registro de metadados de portnum.
10. Preservação da publicação original quando protobuf, validação, descriptografia ou decodificação falham.
11. Política zero-hop opcional via `MESHTASTIC_ZERO_HOP=true`; desativada por padrão.
12. Encerramento mediante SIGINT ou SIGTERM.

## Diferenças deliberadas em relação ao comportamento anterior desta branch

- O handler anterior tentava decodificar toda publicação como `ServiceEnvelope` e retornava `CodeSuccessIgnore` em falhas de parsing ou descriptografia. Isso podia impedir a entrega MQTT de payloads JSON e de pacotes cuja PSK não estivesse disponível localmente. O handler atualizado inspeciona apenas tópicos protobuf reconhecidos e mantém a publicação no fluxo MQTT em caso de falha de inspeção.
- A inspeção local não comprova que um rádio recebeu, descriptografou ou retransmitiu a mensagem. Isso exige testes de interoperabilidade com hardware.
- O certificado PFX e a senha existentes no projeto de referência não foram copiados; certificado e chave privada TLS são configurados por variáveis de ambiente.
- A PSK pode ser configurada por `MESHTASTIC_DEFAULT_PSK` em hexadecimal ou Base64 padrão, desde que resulte em 16 ou 32 bytes. A PSK não é a mesma coisa que a proteção TLS do transporte MQTT.
- `MESHTASTIC_ZERO_HOP` habilita uma transformação opcional de `hop_limit`. A configuração permanece desativada por padrão.

## Segurança do transporte e limites

O listener aceita conexões MQTT sem autenticação adicional. TLS protege o transporte, mas não identifica nem autoriza cada cliente. Restrinja a porta 8883 por firewall/rede confiável e não exponha esta branch diretamente à Internet sem uma camada de autenticação e controles de acesso avaliados separadamente.

Certificados, chaves privadas, arquivos de ambiente e binários locais são excluídos pelo `.gitignore`. Nunca registre nem versione uma PSK de implantação.

A implementação utiliza AES-CTR e uma disposição de nonce validada por teste de interoperabilidade C# registrado em `internal/crypto/crypto_test.go`. Isso não substitui a validação com um pacote real capturado de uma rede Meshtastic e com rádios físicos.

AES-CTR não autentica o conteúdo por si só. Uma descriptografia com PSK incorreta pode produzir bytes que não formam um `Data` protobuf válido. Por isso, a falha de decodificação é tratada como falha de inspeção local, não como evidência conclusiva de corrupção do transporte.

Não há seleção de PSK por `channel_id`; a configuração atual usa uma única PSK local para a tentativa de descriptografia. A PSK local pode não corresponder à chave do canal.

## Tópicos MQTT

O broker não filtra tópicos desconhecidos nem converte JSON automaticamente em um pacote Meshtastic. Publicar um objeto como `{"from":...,"type":"sendtext","payload":"..."}` em `msh/2/json/mqtt/` não cria, por si só, um `ServiceEnvelope` criptografado nem faz downlink para um rádio.

Para observar envelopes publicados no canal `ITUVERAVA`, um assinante pode usar o filtro `msh/2/e/ITUVERAVA/#`. A recepção de uma publicação pelo assinante comprova entrega MQTT, não transmissão LoRa.

## Fora do escopo desta etapa

Ainda não implementados: rate limiting, deduplicação, rate limiting por nó, bloqueio de tópicos desconhecidos, bloqueio de pacotes não descriptografáveis, filtragem de portnums, fail2ban e banimento de atores.

## Critérios de validação pendentes

- `go test ./...`, `go vet ./...` e build no estado atualizado da branch.
- Verificar que payload JSON em `msh/2/json/mqtt/` é encaminhado intacto.
- Verificar que envelopes protobuf inválidos ou não descriptografáveis não são descartados pelo hook.
- Verificar publicação e assinatura nos tópicos do canal de teste.
- Validar downlink e recepção no rádio com hardware real, com PSK e configuração MQTT do dispositivo confirmadas.
