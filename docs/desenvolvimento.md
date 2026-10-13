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
7. Executar go vet ./...
8. Executar go build ./cmd/broker.
9. Rever git diff e git status.
10. Atualizar documentos que descrevem o comportamento alterado.
11. Se houver hardware, registrar modelo, firmware, configuração, procedimento, evidência e limitações em [Equipamentos e ensaios](equipamentos.md) e [Validação](validacao.md).
12. Criar commit com mensagem objetiva.

## Protocolo e criptografia

Mudanças na criptografia ou no processador de pacotes exigem atenção especial. Preservar os testes de interoperabilidade e confirmar alterações no nonce, tamanho da chave, protobuf ou interpretação do pacote. Não presumir suporte a canais com PSKs diferentes sem teste explícito.

## Segurança

A configuração atual não implementa autenticação MQTT nem ACL por tópico. Até que controles adequados sejam implementados e validados, manter o listener em ambiente controlado. TLS não substitui autenticação de clientes nem garante que o cliente valide a identidade do servidor.

## Dependências e documentação

Uma nova dependência precisa de necessidade concreta e licença compatível. Evite duplicar conteúdo: faça links para a fonte canônica da informação. Quando o estado do projeto mudar, revise todos os documentos que dependem dele.

## Critério de pronto

Uma mudança está pronta quando o comportamento esperado está implementado, os testes relevantes foram executados, vet e build foram verificados, a documentação foi atualizada, o diff não contém alterações acidentais e os testes físicos aplicáveis foram realizados ou explicitamente marcados como pendentes.