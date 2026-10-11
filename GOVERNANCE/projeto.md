# Projeto e escopo

## Propósito

O DevNux Mesh Brasil é um projeto de infraestrutura de código aberto para receber e processar publicações Meshtastic por MQTT. O objetivo atual é validar um broker pequeno, compreensível e testável, mantendo evidências claras sobre o que funciona e o que ainda não foi verificado.

## Escopo atual

O código atual implementa um broker MQTT sobre TLS, valida envelopes protobuf, tenta descriptografar pacotes com uma PSK configurada e registra informações dos dados processados. Consulte a documentação técnica para o comportamento exato.

## Fora do escopo confirmado

Autenticação MQTT, ACL por tópico, seleção de PSK por canal, deduplicação, persistência, moderação e rate limiting não devem ser tratados como funcionalidades implementadas. São possibilidades a avaliar, não compromissos automáticos.

## Estado de validação

Há um ensaio físico documentado com um Heltec Wi-Fi LoRa 32 V4. A validação com dois ou mais rádios ainda está pendente. O M5Stack Cardputer Adv está planejado para testes futuros, não registrado como equipamento já validado pelo projeto. Consulte [Equipamentos e ensaios](../docs/equipamentos.md) e [Validação](../docs/validacao.md).

## Autoridade de decisão

Requisitos e mudanças de direção são definidos pelo mantenedor. Sugestões técnicas devem ser identificadas como propostas até serem aceitas. A documentação não substitui a verificação do código, dos testes ou do comportamento físico.