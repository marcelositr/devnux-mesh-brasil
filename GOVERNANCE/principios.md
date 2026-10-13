# Princípios do projeto

## 1. Evidência antes de afirmação

Distinguir implementação presente, teste executado, observação, hipótese e plano futuro. Nunca transformar uma intenção em resultado confirmado.

## 2. Simplicidade deliberada

Adicionar abstrações, dependências e serviços somente quando resolverem uma necessidade demonstrada. Preferir a menor mudança que resolva o problema real.

## 3. Validação incremental

Validar separadamente código, integração MQTT, execução em container e comportamento com hardware. Um tipo de teste não substitui os demais.

## 4. Compatibilidade Meshtastic

Mudanças em protobuf, PSK, nonce, portnums, tópicos ou fluxo MQTT precisam de justificativa técnica e testes de regressão. Não assumir que o broker suporta de forma transparente canais com chaves diferentes.

## 5. Segurança proporcional ao uso

TLS não substitui autenticação MQTT nem autorização por tópico. Até haver controles de acesso adequados e testados, manter o broker em ambiente controlado.

## 6. Histórico preservado

Manter registros de decisões e resultados anteriores. Quando uma conclusão mudar, registrar a nova evidência e explicar a mudança, sem apagar o contexto histórico.

## 7. Responsabilidade humana

Ferramentas automatizadas podem apoiar análise e implementação, mas requisitos e decisões de direção pertencem ao mantenedor.

## 8. Documentação coerente

Quando o comportamento mudar, revisar README, arquitetura, protocolo, configuração, validação, equipamentos e governança que dependam desse comportamento.