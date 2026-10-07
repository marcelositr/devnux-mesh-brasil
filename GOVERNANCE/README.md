# Governança do projeto

Esta pasta contém as regras, princípios e decisões que orientam a evolução do Meshtastic MQTT Brasil.

Ela não é documentação operacional comum. Seu conteúdo define os limites dentro dos quais o projeto deve evoluir.

## Regra de proteção

Nenhuma ferramenta automatizada, incluindo Inteligência Artificial, deve modificar arquivos desta pasta como consequência automática de uma tarefa de implementação.

Mudanças na governança devem ser intencionais, justificadas e registradas separadamente.

## Hierarquia

1. Propósito e princípios do projeto.
2. Regras de governança.
3. Decisões arquiteturais registradas.
4. Documentação técnica.
5. Implementação.

A implementação deve respeitar as decisões vigentes. Quando uma decisão precisar ser alterada, a mudança deve ser documentada antes ou junto da implementação correspondente.

## Conteúdo

- `projeto.md` — finalidade e limites do projeto.
- `principios.md` — princípios permanentes.
- `regras-ia.md` — regras para uso de IA.
- `arquitetura.md` — arquitetura conceitual aprovada.
- `decisoes/` — registro histórico de decisões importantes.

## Regra de preservação

Nenhum documento de governança deve ser apagado apenas porque deixou de refletir a situação atual. Quando uma decisão for substituída, seu registro deve permanecer no histórico e a nova decisão deve explicar a mudança.

> A governança protege a coerência do projeto; a documentação técnica explica como ele funciona; o código implementa o que foi decidido.
