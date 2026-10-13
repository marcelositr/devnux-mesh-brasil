# Governança do DevNux Mesh Brasil

Esta pasta registra o propósito, os princípios e as decisões que orientam a evolução do projeto. A governança não substitui a documentação técnica: ela define como decisões são tomadas e preservadas; a documentação técnica descreve o comportamento implementado; o código deve corresponder ao comportamento verificável.

## Como usar esta pasta

Antes de uma alteração arquitetural ou de segurança:

1. leia [Projeto e escopo](projeto.md);
2. consulte [Princípios](principios.md);
3. procure decisões existentes em [Decisões](decisoes/README.md);
4. descreva o problema e as alternativas;
5. obtenha decisão humana explícita;
6. registre a decisão e atualize os documentos técnicos afetados;
7. valide o resultado de acordo com [Guia de desenvolvimento](../docs/desenvolvimento.md).

## Documentos

- [Projeto e escopo](projeto.md): propósito, público e limites atuais.
- [Princípios](principios.md): critérios duradouros para escolhas técnicas.
- [Arquitetura conceitual](arquitetura.md): responsabilidades e limites da arquitetura.
- [Regras para uso de IA](regras-ia.md): limites e responsabilidades de ferramentas automatizadas.
- [Registro de decisões](decisoes/README.md): índice de decisões arquiteturais e operacionais.

## Hierarquia prática

1. Decisões explícitas do mantenedor e requisitos confirmados.
2. Decisões registradas e ainda vigentes.
3. Princípios do projeto.
4. Documentação técnica atualizada com o comportamento observado.
5. Implementação e testes, que devem ser inspecionados para confirmar o comportamento real.

Se a documentação contradisser o código, não se deve presumir que um deles está correto sem investigação. Registre a divergência e resolva-a de forma explícita.

## Alterações de governança

Mudanças em GOVERNANCE/ devem ser intencionais, revisáveis e justificadas. Ferramentas automatizadas não devem modificar estes documentos como efeito colateral de uma tarefa de implementação. Uma solicitação explícita do mantenedor, como uma revisão dedicada da própria governança, autoriza a revisão documental, mas não autoriza a ferramenta a decidir unilateralmente o rumo do projeto.

Decisões substituídas devem permanecer no histórico; atualize seu estado e registre a decisão sucessora em vez de apagar o registro anterior.

## Regra de evidência

Não registrar uma hipótese como decisão aprovada, nem um teste planejado como teste concluído. Quando faltarem evidências, usar estados explícitos como **proposto**, **pendente**, **observado** ou **validado**, explicando o alcance da afirmação.

> A governança protege a coerência e a rastreabilidade do projeto, não a burocracia pela burocracia.
