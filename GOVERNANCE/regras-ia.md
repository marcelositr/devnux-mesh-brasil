# Regras para uso de Inteligência Artificial

## Regra principal

**A IA deve ajudar a construir o projeto, não decidir o que o projeto deve se tornar.**

Nenhuma IA deve alterar, remover, substituir ou ampliar a arquitetura, os princípios ou o objetivo do projeto por iniciativa própria.

## Proteção da governança

Arquivos dentro de `GOVERNANCE/` não devem ser modificados como consequência automática de uma tarefa de implementação.

Uma alteração de governança exige uma decisão explícita, justificativa e registro.

## Antes de implementar

Quando a tarefa tiver impacto relevante, a IA deve:

1. ler a documentação relacionada;
2. consultar as decisões existentes;
3. identificar restrições aplicáveis;
4. explicar mudanças relevantes;
5. implementar somente o necessário;
6. informar o que foi testado e o que não foi testado.

## Não avacalhar a arquitetura

- Não criar abstrações desnecessárias.
- Não introduzir frameworks ou dependências apenas por conveniência.
- Não reescrever partes funcionais sem necessidade comprovada.
- Não criar funcionalidades sem relação com o objetivo do projeto.
- Não alterar a arquitetura apenas para adequá-la ao padrão preferido da ferramenta.

## Não inventar

- Não afirmar que algo foi testado quando não foi.
- Não inventar resultados de testes.
- Não inventar comportamento de hardware, protocolos ou APIs.
- Em caso de incerteza, declarar a incerteza e propor uma forma de verificar.

## Mudanças incrementais

Preferir alterações pequenas, identificáveis e verificáveis.

Uma tarefa não deve aproveitar uma alteração para reorganizar partes não relacionadas do projeto.

## Histórico e origem

- Não apagar documentação ou decisões históricas para simplificar o projeto.
- Não reescrever o histórico para ocultar a origem.
- Não remover informações de autoria ou licença.
- O uso de IA não transforma a IA em autora do projeto.

## Licença

Todo código incorporado deve respeitar a AGPL-3.0 e ter origem/licença verificável.

Não adicionar código de terceiros sem verificar sua compatibilidade de licença.

## Segurança

Não introduzir credenciais, chaves privadas, tokens ou dados pessoais no repositório.

Não abrir serviços, portas ou integrações externas sem necessidade documentada.

Mudanças com impacto relevante em segurança, disponibilidade ou privacidade devem ser tratadas como mudanças de alto impacto.

## Mudança de decisão

Uma decisão existente pode ser revista quando houver motivo técnico real.

Nesse caso:

1. o problema deve ser identificado;
2. a decisão atual deve ser consultada;
3. a alternativa deve ser explicada;
4. a nova decisão deve ser registrada;
5. somente então a implementação deve seguir.

> **A IA é uma ferramenta de desenvolvimento. A direção do projeto permanece sob decisão humana.**
