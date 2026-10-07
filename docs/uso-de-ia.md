# Regras para uso de IA no projeto

Este projeto pode utilizar ferramentas de Inteligência Artificial como apoio ao desenvolvimento. A IA é uma ferramenta de trabalho e não possui autoridade para definir sozinha o rumo do projeto.

## Regra principal

**Nenhuma IA deve alterar, remover, substituir ou ampliar a arquitetura, os princípios ou o objetivo do projeto por iniciativa própria.**

Antes de uma mudança relevante, a pessoa responsável pelo projeto deve entender a proposta e decidir se ela será incorporada.

## Regras obrigatórias

1. **Não avacalhar a arquitetura**
   - Não criar abstrações desnecessárias.
   - Não introduzir frameworks ou dependências apenas porque são convenientes.
   - Não reescrever partes funcionais sem necessidade comprovada.
   - Não criar funcionalidades que não tenham relação com o objetivo do projeto.

2. **Preservar o que já foi decidido**
   - Consultar README, documentação e decisões anteriores antes de propor alterações.
   - Não contradizer uma decisão documentada sem explicar o motivo.
   - Quando houver conflito entre uma solicitação nova e uma decisão anterior, apontar o conflito antes de alterar o projeto.

3. **Não inventar**
   - Não afirmar que algo foi testado quando não foi.
   - Não inventar resultados de testes, requisitos, APIs, protocolos ou comportamento de hardware.
   - Quando houver incerteza, declarar a incerteza e propor uma forma de verificar.

4. **Mudanças pequenas e verificáveis**
   - Preferir alterações incrementais.
   - Cada mudança deve ter objetivo identificável.
   - Código novo deve ser acompanhado, quando aplicável, de testes ou de uma explicação de como será validado.
   - Evitar mudanças gigantes que misturem várias decisões diferentes.

5. **Não apagar histórico**
   - Não remover documentação, decisões ou autoria para simplificar o projeto.
   - Não reescrever o histórico do projeto para ocultar sua origem.
   - Alterações importantes devem permanecer rastreáveis no Git.

6. **Respeitar a licença**
   - O projeto é distribuído sob AGPL-3.0.
   - Código incorporado ao projeto deve ter licença compatível e origem verificável.
   - Não adicionar código de terceiros sem verificar sua licença e compatibilidade.
   - Não transformar componentes cobertos pela AGPL em componentes proprietários por simples conveniência.

7. **Autoria e origem**
   - O uso de IA não altera a autoria humana registrada no projeto.
   - Ferramentas de IA não devem ser apresentadas como autoras do projeto.
   - Contribuições humanas e decisões relevantes devem continuar identificáveis no histórico quando apropriado.

8. **Segurança e infraestrutura**
   - Não introduzir credenciais, chaves privadas, tokens ou dados pessoais no repositório.
   - Não abrir portas, serviços ou integrações externas sem necessidade documentada.
   - Mudanças que possam afetar disponibilidade, segurança ou privacidade devem ser tratadas como mudanças de alto impacto.

## Antes de implementar

Uma IA trabalhando neste projeto deve, quando a tarefa tiver impacto relevante:

1. ler a documentação relacionada;
2. identificar as decisões existentes;
3. explicar o que pretende mudar;
4. apontar impactos e riscos;
5. implementar somente o necessário;
6. informar o que foi testado e o que não foi testado.

## Princípio do projeto

> **A IA deve ajudar a construir o projeto, não decidir o que o projeto deve se tornar.**

O objetivo é permitir colaboração com pessoas e ferramentas diferentes sem perder coerência, rastreabilidade, simplicidade ou a finalidade comunitária do projeto.
