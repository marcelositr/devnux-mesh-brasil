# Reprodução do Meshtastic MQTT

## Objetivo

Esta etapa reproduz em Go o comportamento efetivamente implementado no projeto de referência meshtastic/mqtt.

O README do projeto de referência também contém ideias futuras que não fazem parte do comportamento atual. Elas não são tratadas como funcionalidades implementadas.

## Comportamentos reproduzidos

1. Broker MQTT.
2. Apenas endpoint MQTT criptografado na porta 8883.
3. TLS 1.2 como protocolo mínimo.
4. Conexões aceitas sem autenticação adicional.
5. Inscrições MQTT aceitas sem filtragem adicional.
6. Bloqueio de payload vazio.
7. Decodificação de ServiceEnvelope em protobuf.
8. Validação dos campos essenciais do envelope.
9. Descriptografia do pacote de canal usando AES-CTR.
10. Decodificação do resultado como Data.
11. Registro de mensagens de texto e de outros portnums.
12. Encerramento mediante SIGINT ou SIGTERM.

## Diferenças deliberadas

- O certificado PFX e a senha existentes no projeto de referência não foram copiados.
- A configuração do certificado foi externalizada por variáveis de ambiente.
- A implementação usa bibliotecas Go em vez das bibliotecas .NET originais.
- A estrutura de código foi separada em pacotes para manter o código testável e idiomático.

## Pontos que precisam de validação

A compatibilidade criptográfica deve ser validada com um pacote real produzido pelo Meshtastic ou com vetores de teste confiáveis.

O projeto de referência usa NonceGenerator da biblioteca Meshtastic para construir o nonce. A implementação Go inicial reproduz a estrutura de 16 bytes usada pelo AES-CTR, mas deve ser confirmada por teste de interoperabilidade antes de ser considerada concluída.

## Fora do escopo atual

As seguintes ideias aparecem no README de referência, mas não estão implementadas no código analisado:

- rate limiting;
- bloqueio de pacotes repetidos;
- rate limiting por nó;
- zero hopping;
- bloqueio de tópicos desconhecidos;
- bloqueio de pacotes não descriptografáveis;
- filtragem de portnums;
- fail2ban;
- banimento de atores.
