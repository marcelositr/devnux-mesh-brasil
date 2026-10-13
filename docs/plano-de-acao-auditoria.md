# Plano de ação — auditoria técnica

- **Projeto:** DevNux Mesh Brasil
- **Branch de trabalho:** `work/finalizacao-projeto`
- **Base da auditoria:** commit `a6f6672c54f83729872e9591b5863f480c0ad650`
- **Estado:** Fase 1 em andamento; comportamento de pacotes sem dados úteis implementado e validado por CI/CodeQL no SHA `11af5e650ca84c615d2d884bc05dfa7b1d8e5c38`.
- **Processo:** trabalhar diretamente na branch de trabalho. Não abrir PRs novos nem fazer merge durante a execução. A integração com `main` será tratada somente no final.
- **Regra de evidência:** não declarar testes executados sem resultado verificável. Separar testes automatizados, testes de container e validação física com hardware.

## Objetivo

Finalizar o software com mudanças pequenas e rastreáveis, comportamento documentado, CI confiável e limitações de segurança explícitas. A validação física multi-nó será feita quando houver rádios disponíveis e não deve ser descrita como concluída antes disso.

## Estado inicial conhecido

- A CI e o CodeQL passaram no commit auditado `a6f6672c54f83729872e9591b5863f480c0ad650`.
- A suíte inclui testes Go, análise estática, build, teste de integração MQTT/TLS com dois clientes e smoke test Docker.
- Há registro de um ensaio físico com um Heltec Wi-Fi LoRa 32 V4.
- A validação com dois ou mais rádios permanece pendente.
- O broker usa TLS, mas ainda não implementa autenticação MQTT nem ACL por tópico; manter em rede controlada.
- O projeto usa AES-CTR, que não autentica criptograficamente o payload.
- O PR #10 já existe, aberto antes da regra de não abrir PRs durante o trabalho. Não fazer merge até a etapa final.

## Ordem de execução

### Fase 1 — Fechar o comportamento do código (primeiro)

- [x] Definir o comportamento: ignorar `Data` com `portnum = 0` ou payload vazio; preservar portnums desconhecidos quando o payload for válido.
- [x] Adicionar testes unitários para descarte de `portnum = 0`/payload vazio e aceitação de portnum desconhecido com payload válido.
- [x] Adicionar verificação de integração MQTT para confirmar que pacote com `portnum = 0` não chega ao assinante.
- [x] Confirmar CI após as alterações: `gofmt`, `go test ./...`, `go vet ./...`, build Go, build Docker e smoke test Docker/TLS concluíram com sucesso na execução 38424940109.
- [x] Preservar o teste que garante que o texto das mensagens não aparece nos logs (`TestLogReceivedMessageDoesNotLogTextPayload`).
- [x] Preservar o vetor de interoperabilidade criptográfica (`TestTransformPacketMatchesInteroperabilityVector`).
- [ ] Revisar o ciclo de vida do servidor e o encerramento por sinais; alterar somente se houver defeito demonstrável.
- [x] Confirmar formatação, `go test ./...`, `go vet ./...` e build via CI.

**Critério de saída:** comportamento definido, coberto por testes pertinentes e coerente com a documentação.

### Fase 2 — Segurança e configuração

- [ ] Avaliar limites de tamanho de pacote, clientes e recursos pendentes com base no tráfego Meshtastic esperado; não escolher valores arbitrários.
- [ ] Se forem definidos limites, configurar e testar os limites no broker.
- [ ] Revisar `.dockerignore` para impedir que arquivos `.env`, certificados, chaves e outros segredos locais entrem no contexto de build, sem excluir exemplos necessários.
- [ ] Documentar permissões e propriedade necessárias para a chave TLS montada no container, sem recomendar permissões excessivas em produção.
- [ ] Manter explícita a restrição a redes confiáveis enquanto não houver autenticação MQTT e ACL.
- [ ] Testar casos de PSK incorreta e payloads inválidos, documentando o que AES-CTR permite ou não detectar.
- [ ] Não alterar a criptografia do protocolo Meshtastic sem requisito e evidência de interoperabilidade.

**Critério de saída:** limites e riscos operacionais explícitos; proteções implementadas têm testes correspondentes.

### Fase 3 — CI, dependências e publicação

- [ ] Revisar o workflow Docker para que a publicação dependa da aprovação das verificações de qualidade.
- [ ] Aplicar o critério às publicações de `latest`, tags de versão e execuções manuais.
- [ ] Verificar vulnerabilidades conhecidas nas dependências Go e nas imagens base.
- [ ] Avaliar a versão da toolchain Go e planejar atualização somente com regressão completa.
- [ ] Executar formatação, testes, vet, build Go, CI, CodeQL e smoke test Docker no commit candidato a release.
- [ ] Revisar o diff integral da branch para confirmar que somente alterações previstas foram incluídas.

**Critério de saída:** commit candidato identificado; CI e CodeQL verdes nesse SHA; imagem construída e testada; publicação condicionada às verificações.

### Fase 4 — Documentação e release

- [ ] Atualizar README, configuração, protocolo, arquitetura, validação, equipamentos e governança quando o comportamento mudar.
- [ ] Registrar resultados efetivamente obtidos e itens que continuam pendentes.
- [ ] Definir os critérios da primeira versão e as instruções de implantação em rede controlada.
- [ ] Definir um canal de reporte de vulnerabilidades antes de uma distribuição pública mais ampla.
- [ ] Revisar o PR #10 e decidir a integração com `main` somente depois da conclusão das fases anteriores.

**Critério de saída:** documentação coerente, evidências rastreáveis e revisão final explícita da integração.

### Fase 5 — Validação física posterior

- [ ] Repetir o ensaio com o Heltec Wi-Fi LoRa 32 V4 e registrar estabilidade, reconexões e resultados.
- [ ] Quando houver segundo rádio, validar o fluxo multi-nó, incluindo uplink e downlink conforme o cenário suportado.
- [ ] Investigar mensagens duplicadas com IDs, timestamps e campos dos pacotes; não implementar deduplicação sem evidência.
- [ ] Atualizar o registro de equipamentos e a documentação de validação.

**Critério de saída:** resultados observados com dois ou mais nós, sem declarar antecipadamente a compatibilidade multi-nó.

## Achados priorizados

### P1 — antes de considerar a finalização

1. **Publicação Docker separada da CI:** estabelecer uma barreira para não publicar uma imagem se as verificações de qualidade falharem.
2. **Ausência de autenticação MQTT e ACL:** manter a rede controlada; decidir os requisitos de identidade e autorização antes de qualquer exposição pública.

### P2 — resolver ou decidir antes do release

3. **Limites de recursos:** avaliar limites de pacote, clientes e recursos pendentes com base no tráfego pretendido.
4. **Pacotes sem dados úteis — resolvido na Fase 1:** são ignorados; portnums desconhecidos com payload válido continuam sendo encaminhados.
5. **Contexto Docker:** revisar exclusões de arquivos de ambiente e segredos.
6. **Permissões da chave TLS:** documentar um procedimento compatível com usuário não privilegiado no container.
7. **AES-CTR sem autenticação:** testar falhas de processamento e documentar limites sem quebrar o protocolo Meshtastic.

### P3 — melhorias adicionais

8. Verificar suporte e atualização da toolchain Go.
9. Adicionar verificação de vulnerabilidades de dependências e imagens base.
10. Considerar `go test -race` e fuzzing direcionado a parsers como testes complementares.
11. Avaliar fixação de imagens por digest e geração de SBOM para release.
12. Definir canal de reporte de segurança.

## Decisões humanas pendentes

Antes de implementar decisões que alterem o escopo, confirmar:

1. **Ambiente de uso:** somente rede local/controlada ou futura exposição pública?
2. **Semântica de publicação:** pacotes sem `Data` útil devem ser descartados ou permanecer no fluxo MQTT?
3. **Limites operacionais:** qual é o tamanho/carga de pacotes esperado para o cenário de uso?

O trabalho de auditoria, testes automatizados e correções não deve ser confundido com validação física multi-nó.

## Registro de progresso

Atualizar esta seção a cada etapa concluída, com o SHA do commit e links dos resultados:

- **Fase atual:** Fase 1 — comportamento e documentação atualizados; CI e CodeQL verdes. Falta revisar ciclo de vida/encerramento do servidor antes de fechar a fase.
- **SHA inicial da auditoria:** `a6f6672c54f83729872e9591b5863f480c0ad650`.
- **SHA com testes de regressão antes da implementação:** `8fd8aa9d823693243d29ffb100adfbcc5cc901e1` — CI falhou como esperado em `TestOnPublishIgnoresDataWithoutUsefulPayload`, confirmando o defeito.
- **SHA candidato validado:** `11af5e650ca84c615d2d884bc05dfa7b1d8e5c38`.
- **CI:** [execução 38424940109](https://github.com/marcelositr/devnux-mesh-brasil/actions/runs/38424940109) — sucesso; job `Testes e qualidade` concluído.
- **CodeQL:** [execução 38424940149](https://github.com/marcelositr/devnux-mesh-brasil/actions/runs/38424940149) — sucesso; job `Análise CodeQL` concluído.
- **Revisão do diff:** comparação com o SHA inicial identificou apenas os seis arquivos previstos para testes, comportamento e documentação.
- **Próxima ação:** revisar o ciclo de vida do servidor e o encerramento por sinais; depois reavaliar o critério de saída da Fase 1.
- **Validação física multi-nó:** pendente.
- **Integração com `main`:** reservada para o final.
