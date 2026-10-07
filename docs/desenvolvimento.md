# Guia de desenvolvimento

## Objetivo

Este documento define como modificar e validar o DevNux Mesh Brasil sem perder a simplicidade da base atual.

A arquitetura está em [Arquitetura](arquitetura.md). O comportamento do protocolo está em [Protocolo e criptografia](protocolo.md). A execução está em [Configuração e execução](configuracao.md).

## Princípios

### Preservar comportamento

Alterações internas devem manter o comportamento funcional existente, salvo quando a tarefa pedir explicitamente uma mudança de comportamento.

### Preferir mudanças pequenas

Uma alteração deve ter um objetivo identificável. Evite misturar refatoração, mudança de protocolo e nova funcionalidade no mesmo conjunto de mudanças.

### Não abstrair por antecipação

Interfaces, fábricas, camadas e novos pacotes devem existir quando houver uma necessidade concreta. Não introduza uma abstração apenas para tornar a arquitetura mais genérica.

### Testar antes de considerar concluído

Código alterado deve ser formatado e validado.

## Fluxo recomendado

1. Identificar a responsabilidade afetada.
2. Ler a documentação correspondente.
3. Alterar somente o necessário.
4. Atualizar ou adicionar testes quando o comportamento exigir.
5. Executar `gofmt`.
6. Executar `go test ./...`.
7. Executar `go vet ./...`.
8. Executar `go build ./cmd/broker`.
9. Conferir `git diff` e `git status`.
10. Registrar no commit uma descrição objetiva da mudança.

## Comandos de validação

```bash
gofmt -w cmd/broker internal
go test ./...
go vet ./...
go build -o /tmp/devnux-mesh-broker ./cmd/broker
```

## Alterações de protocolo e criptografia

Mudanças em `internal/crypto` ou `internal/broker/packet_processor.go` exigem atenção especial.

Além da bateria geral, preserve os testes de interoperabilidade criptográfica e confirme explicitamente qualquer alteração no formato do nonce, tamanho de chave, protobuf ou interpretação do pacote.

Consulte [Protocolo e criptografia](protocolo.md) antes de alterar essas áreas.

## Dependências

Uma nova dependência deve ter necessidade concreta e ser compatível com o objetivo do projeto.

Evite dependências usadas apenas para substituir poucas linhas de código ou introduzir uma abstração que o projeto não necessita.

## Documentação

A documentação deve explicar decisões e comportamentos que não sejam óbvios pelo código.

Evite repetir o conteúdo de outro documento. Quando uma informação já estiver documentada, faça um link para ela.

## Histórico

Não reescreva o histórico Git para ocultar decisões ou mudanças. Commits devem permitir compreender a evolução técnica do projeto.

## Critério de pronto

Uma mudança está pronta quando:

- o comportamento esperado está implementado;
- os testes relevantes passam;
- `go vet` passa;
- o broker compila;
- a documentação necessária foi atualizada;
- não existem alterações acidentais no diff.
