# Evolução futura e hardening

## Objetivo

Registrar possibilidades de evolução do broker sem transformar hipóteses em requisitos de implementação.

O DevNux Mesh Brasil deve primeiro comprovar, com hardware Meshtastic físico e real, que o comportamento atual do broker corresponde ao funcionamento esperado da rede. Somente depois dessa validação será feita a curadoria das medidas de hardening e moderação que realmente forem necessárias.

## Regra de evolução

As funcionalidades deste documento são **possibilidades futuras**. Elas não fazem parte da implementação atual e não devem ser implementadas apenas porque aparecem nesta lista.

O fluxo de decisão é:

```text
broker atual
    ↓
testes com hardware Meshtastic real
    ↓
observação do tráfego e do comportamento reais
    ↓
identificação de problemas concretos
    ↓
curadoria do que realmente é necessário
    ↓
implementação incremental
    ↓
nova validação com hardware real
```

A razão é prática: implementar mecanismos de controle antes de observar o comportamento real pode criar regras incorretas e gerar retrabalho quando o hardware chegar.

## Possibilidades futuras

A lista abaixo é deliberadamente curta e priorizada. Ela não representa um compromisso de implementar todos os itens.

### Prioridade provável

#### Deduplicação

Investigar e, se necessário, implementar deduplicação de pacotes para evitar processamento ou redistribuição repetida do mesmo tráfego.

A regra de deduplicação, a identidade usada para reconhecer um pacote e o tempo de retenção devem ser definidos somente depois de observar o comportamento real da rede.

#### Rate limiting

Avaliar limites de tráfego por cliente, nó ou outra unidade que faça sentido após os testes reais.

Os limites não devem ser definidos por valores arbitrários antes de conhecer o comportamento normal da rede. O objetivo é reduzir abuso sem bloquear tráfego legítimo.

#### Observabilidade

Evoluir os registros atuais para permitir identificar comportamento anômalo e operar o broker com segurança.

Antes de adicionar mecanismos automáticos de bloqueio, deve ser possível observar pelo menos origem, cliente, tópico e volume de eventos relevantes, conforme o que o hardware e o protocolo demonstrarem ser necessário.

#### ACL e autorização

Avaliar controle de publicação e inscrição por cliente e tópico antes de uma exposição pública mais ampla.

A política deve ser simples e baseada no uso real do broker, evitando criar uma camada de autorização maior do que o problema exige.

### Possibilidades condicionais

#### Bloqueio ou banimento de clientes

Pode ser necessário caso os testes e a operação real revelem clientes persistentemente abusivos.

Qualquer mecanismo de bloqueio deve nascer de evidência observável e possuir critérios claros. Não deve ser implementado preventivamente.

#### Moderação

Mecanismos de moderação podem ser considerados se houver necessidade operacional real.

Moderação não deve ser confundida com filtragem arbitrária de pacotes legítimos.

#### Persistência

Pode ser considerada caso surja uma necessidade concreta de histórico, auditoria, recuperação ou outra função que justifique manter estado fora da memória.

Não é requisito do broker atual.

### Baixa prioridade ou não planejado

#### Filtragem avançada de portnums

Não é objetivo atual criar uma lista extensa de `portnum` permitidos ou proibidos.

Isso só deve ser reconsiderado se os testes com hardware demonstrarem um problema concreto que não possa ser resolvido de forma mais simples.

#### Controle de hop

O broker não deve, neste estágio, tentar reproduzir no servidor regras de encaminhamento próprias do rádio para decidir quais pacotes devem ou não fazer hop.

Esse comportamento pertence ao funcionamento da rede Meshtastic e só deveria ser tratado no broker se houver requisito técnico concreto e claramente delimitado.

#### Fail2ban ou mecanismos equivalentes

Não há motivo para introduzir mecanismos desse tipo antes de existir evidência de abuso que justifique a complexidade operacional.

## O que já existe

A implementação atual já rejeita payloads vazios, protobufs inválidos, envelopes inválidos, pacotes sem payload criptografado e pacotes que não podem ser interpretados pelo processamento atual.

Essas validações fazem parte do funcionamento atual e não devem ser confundidas com o hardening futuro descrito aqui.

## Critério para sair desta fase

O projeto só deve passar da fase de observação para a implementação de hardening depois que:

- o broker atual estiver validado com hardware Meshtastic físico e real;
- o fluxo MQTT real estiver observado;
- o comportamento normal dos nós e gateways estiver compreendido;
- problemas concretos tiverem sido identificados;
- cada nova medida tiver uma justificativa baseada nesses dados;
- a medida puder ser testada sem comprometer a compatibilidade já validada.

A lista deste documento deve ser revisada nesse momento. Itens que não demonstrarem necessidade devem ser descartados, mesmo que estejam aqui registrados.

## Princípio

> Primeiro provar que o broker atual funciona com hardware real. Depois decidir o que precisa ser protegido, limitado ou moderado.
