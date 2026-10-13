# Guia de desenvolvimento

## Objetivo

Este documento define como modificar e validar o DevNux Mesh Brasil sem perder a simplicidade da base atual. Consulte [Arquitetura](arquitetura.md), [Protocolo e criptografia](protocolo.md), [Configuração](configuracao.md), [Validação](validacao.md) e [Equipamentos e ensaios](equipamentos.md). As decisões e princípios do projeto ficam em [GOVERNANCE](../GOVERNANCE/README.md).

## Princípios

### Preservar comportamento

Alterações internas devem manter o comportamento funcional existente, salvo quando a tarefa pedir explicitamente uma mudança.

### Preferir mudanças pequenas

Cada alteração deve ter objetivo identificável. Evite misturar refatoração, mudança de protocolo e nova funcionalidade no mesmo conjunto de mudanças.

### Não abstrair por antecipação

Interfaces, fábricas, camadas e novos pacotes devem existir quando houver necessidade concreta.

### Distinguir tipos de validação

Testes unitários, testes de integração MQTT, testes em Docker e ensaios com hardware físico são evidências diferentes. Um teste bem-sucedido com um rádio não comprova o encaminhamento em uma malha com vários nós.

## Fluxo recomendado

1. Identificar a responsabilidade afetada.
2. Ler a documentação e as decisões aplicáveis.
3. Alterar somente o necessário.
4. Atualizar ou adicionar testes quando o comportamento exigir.
5. Executar gofmt.
6. Executar go test ./...
7. Executar `go vet ./...`.
8. Executar `go test -race ./...` para procurar condições de corrida.
9. Executar `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` e analisar vulnerabilidades conhecidas que afetem funções chamadas.
10. Executar `go build ./cmd/broker`.
11. Validar o contexto Docker, construir a imagem, executar o scanner de vulnerabilidades e testar a inicialização TLS.
12. Rever o diff integral e o estado do Git; não incluir scripts ou artefatos locais sem aprovação.
13. Atualizar documentos que descrevem o comportamento alterado.
14. Se houver hardware, registrar modelo, firmware, configuração, procedimento, evidência e limitações em [Equipamentos e ensaios](equipamentos.md) e [Validação](validacao.md).
15. Criar commit com mensagem objetiva.

## Protocolo e criptografia

Mudanças na criptografia ou no processador de pacotes exigem atenção especial. Preservar os testes de interoperabilidade e confirmar alterações no nonce, tamanho da chave, protobuf ou interpretação do pacote. Não presumir suporte a canais com PSKs diferentes sem teste explícito.

## Segurança

A configuração atual não implementa autenticação MQTT nem ACL por tópico. Até que controles adequados sejam implementados e validados, manter o listener em ambiente controlado. TLS não substitui autenticação de clientes nem garante que o cliente valide a identidade do servidor.

## Dependências e documentação

Uma nova dependência precisa de necessidade concreta e licença compatível. Evite duplicar conteúdo: faça links para a fonte canônica da informação. Quando o estado do projeto mudar, revise todos os documentos que dependem dele.

## Critério de pronto

Uma mudança está pronta quando o comportamento esperado está implementado, testes, `vet`, `race`, verificação de vulnerabilidades e build foram verificados, a CI e o CodeQL do commit candidato estão verdes, a imagem Docker e seu scan foram verificados, a documentação está atualizada e o diff não contém alterações acidentais. Os testes físicos aplicáveis devem ser realizados ou explicitamente marcados como pendentes.