# Plano de ação — auditoria técnica

- **Projeto:** DevNux Mesh Brasil
- **Branch de trabalho:** `work/finalizacao-projeto`
- **Base da auditoria:** commit `a6f6672c54f83729872e9591b5863f480c0ad650`
- **Estado:** verificações automatizáveis das Fases 2 e 3 concluídas e aprovadas em CI/CodeQL/Fuzz no candidato `8af16e8`; documentação e critérios de release revisados. Restam uma configuração manual de reporte privado no GitHub e a decisão explícita sobre merge do PR #10. A Fase 5 (hardware e carga real) fica para depois. Meta de planejamento: até 200 rádios, sem tratar isso como capacidade comprovada nem converter diretamente rádios em conexões MQTT.
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
- [x] Revisar o ciclo de vida do servidor e o encerramento por sinais; `main.go` trata falha de inicialização, registra `SIGINT`/`SIGTERM` e chama `server.Close()`. Não foi identificado defeito demonstrável que justifique alterar a implementação nesta etapa.
- [x] Confirmar formatação, `go test ./...`, `go vet ./...` e build via CI.

**Critério de saída:** comportamento definido, coberto por testes pertinentes e coerente com a documentação.

### Fase 2 — Segurança e configuração

- [x] Consultar a documentação oficial do Mochi MQTT v2.7.9: `MaximumPacketSize: 0` não limita o pacote MQTT; `MaximumClientWritesPending` tem padrão 8192 por cliente; `MaximumClients` não impõe um teto operacional útil. Registrado em `docs/configuracao.md`.
- [x] Documentar a metodologia para coletar tráfego real e testar conexões, carga normal, pico observado, assinantes lentos e reconexões em `docs/perfil-de-carga.md`.
- [x] Revisar `.dockerignore` para excluir arquivos `.env`, certificados, chaves e diretórios locais de segredos, preservando `.env.example`; o teste de regressão `scripts/test-dockerignore.sh` passou na CI do candidato `8af16e8`.
- [x] Documentar a recomendação de menor privilégio para a chave TLS montada: chave `0640` restrita a root/grupo dedicado, acesso do grupo no container e bind mount somente leitura.
- [x] Confirmar CI do smoke test Docker com chave `0640` e grupo correspondente ao usuário `devnux` — sucesso na execução [38429481414](https://github.com/marcelositr/devnux-mesh-brasil/actions/runs/38429481414).
- [x] Manter explícita a restrição a redes confiáveis enquanto não houver autenticação MQTT e ACL; consta no README, configuração, desenvolvimento e SECURITY.md.
- [x] Adicionar teste determinístico de PSK incorreta e documentar payloads inválidos e os limites de AES-CTR. O teste não afirma detecção criptográfica geral; CI passou no commit candidato.
Não alterar a criptografia do protocolo Meshtastic sem requisito e evidência de interoperabilidade.

**Critério de saída:** limites e riscos operacionais explícitos; proteções implementadas têm testes correspondentes.

### Fase 3 — CI, dependências e publicação

- [x] Revisar o workflow Docker para que a publicação dependa da aprovação das verificações de qualidade.
- [x] Aplicar o gate às publicações de `latest`, tags de versão e execuções manuais; publicação manual fora de `main` é ignorada.
- [x] Verificar vulnerabilidades conhecidas nas dependências Go e na imagem Docker. No commit candidato, `govulncheck` não encontrou vulnerabilidades alcançáveis e o scan Trivy passou para os níveis configurados (CRITICAL/HIGH com correção disponível).
- [x] Atualizar a toolchain Go para 1.27.2 após a identificação de vulnerabilidades alcançáveis na biblioteca padrão 1.27.1; a regressão completa passou.
- [x] Fixar as Actions de terceiros por SHA completo nos workflows CI, CodeQL, Docker e Fuzz, além da action de upload de artefatos; os workflows passaram após a alteração.
- [x] Fixar por digest SHA-256 as imagens-base Go 1.27.2-alpine e Alpine 3.22; o build, scan e smoke test passaram.
- [x] Gerar SBOM CycloneDX da imagem Docker e publicar como artefato de CI com retenção de 30 dias.
- [ ] Habilitar/confirmar o reporte privado de vulnerabilidades nas configurações do repositório; a automação disponível não tem permissão para alterar esta configuração.
- [x] Executar formatação, testes, race, govulncheck, vet, build Go, CI, CodeQL, scan, SBOM e smoke test Docker no candidato `8af16e89ec2b522fcf887f5b836b6eace47d1fcc`.
- [x] Revisar a lista integral de arquivos alterados entre `main` e `work/finalizacao-projeto`; os 23 arquivos alterados correspondem às correções, testes, workflows, segurança e documentação previstas no plano. A validação final do diff permanece ligada ao commit candidato e não inclui scripts locais não rastreados.

**Critério de saída:** commit candidato identificado; CI e CodeQL verdes nesse SHA; imagem construída e testada; publicação condicionada às verificações.

### Fase 4 — Documentação e release

- [x] Atualizar README, configuração, protocolo, arquitetura e validação para refletir comportamento, toolchain e cobertura de testes; equipamentos e governança foram revisados e permanecem coerentes com a fase física pendente.
- [x] Registrar resultados efetivamente obtidos e itens que continuam pendentes, com links de CI/CodeQL.
- [x] Documentar critérios técnicos de conclusão e instruções de implantação em rede controlada.
- [ ] Habilitar o reporte privado de vulnerabilidades nas configurações do GitHub; a documentação e o link direto estão prontos, mas a ativação exige ação do proprietário do repositório.
- [x] Revisar o PR #10; decisão de merge permanece pendente de autorização explícita do usuário.

**Critério de saída:** documentação coerente, evidências rastreáveis e revisão final explícita da integração.

### Fase 5 — Validação física posterior

- [ ] Definir limites de tamanho de pacote, clientes e recursos pendentes após medir envelopes reais, frequência de publicação e consumo sob carga para a meta de planejamento de até 200 rádios; não escolher valores arbitrários.
- [ ] Definir o modelo de conexão da rede (rádios conectados diretamente ou por gateways), pois o número de rádios não equivale necessariamente ao número de clientes MQTT.
- [ ] Coletar dados reais de tamanho/frequência dos envelopes e medir memória/CPU/latência antes de declarar capacidade suportada.
- [ ] Se forem definidos limites, configurar e testar os limites no broker.
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
11. [x] Fixar imagens-base por digest SHA-256 e gerar SBOM CycloneDX como artefato de CI.
12. Definir canal de reporte de segurança.

## Decisões humanas pendentes

Antes de implementar decisões que alterem o escopo, confirmar:

1. **Ambiente de uso:** somente rede local/controlada ou futura exposição pública?
2. **Semântica de publicação:** pacotes sem `Data` útil devem ser descartados ou permanecer no fluxo MQTT?
3. **Meta inicial de escala — definida para planejamento:** até 200 rádios. A capacidade real e os limites operacionais ainda precisam ser validados por perfil de carga e testes; a topologia MQTT/gateways precisa ser considerada.

O trabalho de auditoria, testes automatizados e correções não deve ser confundido com validação física multi-nó.

## Registro de progresso

Atualizar esta seção a cada etapa concluída, com o SHA do commit e links dos resultados:

- **Fase atual:** trabalho automatizável das Fases 2, 3 e 4 concluído e verificado no candidato `8af16e8`. A Fase 5 (medições de carga e validação física multi-rádio) permanece explicitamente posterior.
- **SHA inicial da auditoria:** `a6f6672c54f83729872e9591b5863f480c0ad650`.
- **SHA com testes de regressão antes da implementação:** `8fd8aa9d823693243d29ffb100adfbcc5cc901e1` — CI falhou como esperado em `TestOnPublishIgnoresDataWithoutUsefulPayload`, confirmando o defeito.
- **SHA candidato validado:** `11af5e650ca84c615d2d884bc05dfa7b1d8e5c38`.
- **CI:** [execução 38424940109](https://github.com/marcelositr/devnux-mesh-brasil/actions/runs/38424940109) — sucesso; job `Testes e qualidade` concluído.
- **CodeQL:** [execução 38424940149](https://github.com/marcelositr/devnux-mesh-brasil/actions/runs/38424940149) — sucesso; job `Análise CodeQL` concluído.
- **Revisão do diff:** comparação com o SHA inicial identificou apenas os seis arquivos previstos para testes, comportamento e documentação.
- **Revisão do ciclo de vida:** `cmd/broker/main.go` já lida com falha de inicialização e encerramento por `SIGINT`/`SIGTERM`; nenhum defeito demonstrável encontrado, sem alteração de código.
- **Decisão registrada:** planejar inicialmente para uma rede de até 200 rádios. Não definir limites de pacote/conexões até examinar a documentação da versão Mochi MQTT usada e estabelecer perfil de carga plausível.
- **Progresso Fase 2:** teste de PSK incorreta adicionado com ressalva explícita de que AES-CTR não autentica conteúdo; restrição de rede controlada documentada; CI atual passou no teste de contexto Docker e no smoke test TLS com chave `0640`. Os defaults do Mochi MQTT v2.7.9 estão documentados; não foram introduzidos limites numéricos arbitrários.
- **Progresso de dimensionamento:** a metodologia de teste de carga foi registrada em `docs/perfil-de-carga.md`; ainda faltam topologia final, captura de tráfego representativa e medições do host.
- **Próximas ações:** habilitar/confirmar o reporte privado de vulnerabilidades nas configurações do GitHub e decidir explicitamente se o PR #10 deve ser integrado. Perfil de carga e limites numéricos foram movidos para a Fase 5, pois dependem de tráfego real e medições. Não alterar a cifra Meshtastic.
- **Candidato automatizado validado:** `8af16e89ec2b522fcf887f5b836b6eace47d1fcc` — [CI](https://github.com/marcelositr/devnux-mesh-brasil/actions/runs/38445485210), [CodeQL](https://github.com/marcelositr/devnux-mesh-brasil/actions/runs/38445492503) e [Fuzz](https://github.com/marcelositr/devnux-mesh-brasil/actions/runs/38445485274) passaram. A CI gerou o SBOM CycloneDX como artefato.
- **Validação física multi-nó:** pendente.
- **PR #10:** revisado; não foi feito merge. A integração fica pendente de autorização explícita.
- **Reporte privado de vulnerabilidades:** ainda requer ativação manual em [Settings → Security](https://github.com/marcelositr/devnux-mesh-brasil/settings/security_analysis).
