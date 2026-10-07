# Meshtastic MQTT Brasil

Infraestrutura MQTT comunitária para integração e comunicação entre redes Meshtastic no Brasil, com foco em uso público, colaboração e continuidade da comunidade.

## Status

**Em planejamento.**

Este repositório é o ponto de partida para documentar, desenvolver e acompanhar uma futura infraestrutura MQTT voltada ao ecossistema Meshtastic no Brasil.

## Objetivo

Estudar e, futuramente, desenvolver uma infraestrutura que permita interligar diferentes redes e regiões Meshtastic por meio de MQTT, sem substituir o funcionamento local das redes LoRa.

A ideia é que a comunicação local continue funcionando mesmo quando a conexão com a infraestrutura MQTT ou com a Internet estiver indisponível.

## Princípios

- Projeto voltado à comunidade.
- Código e conhecimento devem permanecer acessíveis à comunidade conforme os termos da licença.
- A infraestrutura deve complementar as redes locais, não depender delas para funcionar.
- Prioridade para continuidade, simplicidade e funcionamento em situações reais.
- Desenvolvimento colaborativo e documentado.
- Respeito à origem e ao histórico do projeto.

## Possível arquitetura

```text
              NÚCLEO MQTT
                    │
        ┌───────────┴───────────┐
        │                       │
     Região A                Região B
        │                       │
     Gateway                 Gateway
        │                       │
       LoRa                    LoRa
        │                       │
    Rede local             Rede local
        │                       │
    Dispositivos            Dispositivos
```

A arquitetura definitiva ainda não está definida. Ela deverá ser construída a partir de testes com hardware real, redes Meshtastic reais e gateways funcionando em condições práticas.

## Desenvolvimento

O desenvolvimento deverá seguir uma evolução gradual:

1. Conhecer e testar o Meshtastic em hardware real.
2. Construir e observar redes locais.
3. Testar gateways.
4. Estudar a integração com MQTT.
5. Testar a comunicação entre regiões.
6. Avaliar disponibilidade, segurança e continuidade.
7. Documentar os resultados.
8. Somente então considerar uma infraestrutura de maior escala.

## Situação atual

Nenhuma infraestrutura nacional está sendo implantada neste momento.

Este repositório existe para preservar a ideia, organizar o planejamento e servir como base para um possível desenvolvimento futuro.

## Licença

Este projeto utiliza a **GNU Affero General Public License v3.0 (AGPL-3.0)**.

A licença permite uso, estudo, modificação e distribuição do projeto, inclusive para fins comerciais, desde que sejam respeitadas as condições da AGPL-3.0, incluindo as obrigações relacionadas à disponibilização do código-fonte correspondente em versões modificadas oferecidas como serviço pela rede.

Consulte o arquivo [`LICENSE`](LICENSE) para os termos completos.

## Autoria e origem

Projeto iniciado por **marcelositr**.

O histórico público deste repositório faz parte da documentação da origem e da evolução do projeto. Contribuições posteriores devem respeitar a licença e preservar os avisos de autoria e licença aplicáveis.

Eventuais registros formais de versões do software junto ao INPI poderão ser realizados conforme a evolução do projeto.

## Contribuições

Contribuições, testes, documentação, sugestões e implementações poderão ser incorporados ao projeto conforme sua evolução.

Antes de iniciar uma implementação de grande porte, consulte as discussões e a documentação existentes para evitar trabalho duplicado e manter uma arquitetura coerente.

## Próximos passos

- [ ] Definir requisitos iniciais.
- [ ] Testar hardware Meshtastic.
- [ ] Testar uma rede local.
- [ ] Estudar o papel do gateway.
- [ ] Definir a primeira arquitetura MQTT.
- [ ] Criar documentação técnica.
- [ ] Realizar testes de disponibilidade e recuperação.
- [ ] Avaliar expansão para outras regiões.

---

**Meshtastic MQTT Brasil**  
Infraestrutura comunitária em planejamento para integração de redes Meshtastic no Brasil.
