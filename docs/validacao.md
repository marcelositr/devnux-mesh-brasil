# Validação

Este documento registra as validações funcionais realizadas no broker durante o desenvolvimento.

## Validação local

A execução direta no host foi validada com:

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

Esse teste confirmou o fluxo:

```text
cliente MQTT
    -> TLS
    -> broker no container
    -> ServiceEnvelope
    -> descriptografia Meshtastic
    -> processamento
    -> publicação MQTT
```

Após o teste, a permissão da chave privada temporária deve ser restaurada para `0600`.

## Limites da validação

A validação confirma o funcionamento da imagem Docker e do fluxo principal do broker. Ela não substitui testes de implantação em ambiente público.

Continuam fora do escopo atual mecanismos como autenticação MQTT, ACL, rate limiting, deduplicação, persistência, moderação e observabilidade avançada.
