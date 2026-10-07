# Comportamento e protocolo MQTT do DevNux Mesh Brasil

## Objetivo

Documentar o comportamento atualmente implementado pelo DevNux Mesh Brasil e os limites técnicos desta etapa.

## Comportamentos implementados

1. Broker MQTT.
2. Endpoint MQTT sobre TLS na porta 8883.
3. TLS 1.2 como protocolo mínimo.
4. Conexões aceitas sem autenticação adicional.
5. Inscrições MQTT aceitas sem filtragem adicional.
6. Payload vazio ignorado.
7. Decodificação de `ServiceEnvelope` em protobuf.
8. Validação dos campos essenciais do envelope.
9. Descriptografia do pacote usando AES-CTR.
10. Decodificação do resultado como `Data`.
11. Registro de mensagens de texto e de outros portnums.
12. Encerramento mediante SIGINT ou SIGTERM.

## Configuração de certificados

O certificado e a chave privada são arquivos separados e configuráveis por ambiente:

- `MQTT_CERTIFICATE_FILE`
- `MQTT_PRIVATE_KEY_FILE`

A configuração padrão utiliza `certificate.pem` e `private.key`.

## Criptografia

A implementação utiliza AES-CTR.

O nonce possui 16 bytes e é construído a partir do identificador do pacote e da origem, em little-endian, mantendo os quatro bytes finais zerados.

A PSK padrão possui 16 bytes:

`d4f1bb3a20290759f0bcffabcf4e6901`

Também é possível configurar uma PSK de 16 ou 32 bytes pela variável `MESHTASTIC_DEFAULT_PSK`, em hexadecimal.

A implementação possui teste de interoperabilidade para validar a construção do nonce e a descriptografia de um vetor conhecido.

## Validação de pacotes

Antes da descriptografia, o broker verifica:

- `channel_id` preenchido;
- `gateway_id` preenchido;
- pacote presente;
- identificador do pacote válido;
- origem válida;
- payload criptografado presente;
- ausência de payload decodificado no pacote.

Após a descriptografia, o conteúdo é interpretado como `Data`. Dados sem portnum válido ou sem payload são aceitos pelo processamento, mas não geram registro como mensagem de texto.

## Tópicos MQTT

O nome do tópico não possui filtragem específica no processamento das publicações. Ele é preservado e utilizado no registro da mensagem.

Não há seleção de PSK por `channel_id`. A descriptografia utiliza a PSK padrão configurada para o broker.

## Limites atuais

Os seguintes recursos ainda não fazem parte da implementação:

- rate limiting;
- deduplicação de pacotes;
- persistência;
- ACL ou autorização por tópico;
- filtragem de portnums;
- observabilidade avançada;
- mecanismos de banimento ou moderação.

Esses recursos podem ser tratados em etapas futuras sem alterar o núcleo de processamento já validado.
