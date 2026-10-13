# Evolução futura e hardening

## Objetivo

Registrar possibilidades de evolução sem transformar hipóteses em requisitos automáticos. A decisão atual é concluir e verificar o software com os testes automatizados, de integração, build e Docker que podem ser executados sem comprar ou reunir mais rádios. A validação física com dois ou mais nós fica registrada como etapa posterior, não como bloqueio para esse trabalho. Consulte [Validação](validacao.md) e [Equipamentos e ensaios](equipamentos.md).

## Correções técnicas identificadas na revisão estática

Os itens abaixo são observações de leitura do código, não resultados de testes de execução. Devem ser confirmados com testes antes de declarar uma correção concluída:

- **Erro ao iniciar o servidor:** o retorno de Serve é registrado numa goroutine, mas o processo principal continua aguardando SIGINT/SIGTERM. Se Serve terminar com erro, o processo pode permanecer ativo sem serviço funcional. Tratar o erro no ciclo de vida e cobrir o comportamento com teste.
- **Logs de mensagens:** o hook registra o conteúdo integral das mensagens de texto em nível INFO. Remover o conteúdo da mensagem dos logs por padrão para reduzir exposição de conversas e crescimento dos logs.
- **Teste de integração MQTT:** os testes do hook não substituem um cenário com dois clientes MQTT conectados, um publicando e outro recebendo a publicação. Cobrir esse fluxo em teste automatizado.
- **Limites de entrada:** verificar e documentar os limites efetivos de tamanho de publicação e conexões na biblioteca MQTT utilizada; não presumir que existam limites específicos configurados pelo projeto.

## Sequência de finalização do software

1. Corrigir afirmações documentais que excedem a cobertura real dos testes.
2. Tratar o ciclo de vida do servidor quando `Serve` termina com erro e adicionar cobertura automatizada.
3. Evitar registrar o conteúdo integral de mensagens de texto em logs operacionais.
4. Criar teste de integração com dois clientes MQTT para comprovar publicação e recebimento.
5. Verificar limites relevantes de entrada e conexões da biblioteca usada e documentar o que estiver efetivamente configurado.
6. Executar formatação, testes, `go vet`, build, CI, CodeQL e construção/execução Docker conforme os recursos disponíveis.
7. Revisar os documentos, exemplos de configuração, segurança, publicação da imagem e critérios de release.
8. Registrar explicitamente os testes físicos ainda pendentes e os limites de uso da versão finalizada.

## Regra de evolução

As funcionalidades listadas abaixo são possibilidades, não compromissos. Não implementar um item apenas porque aparece neste documento. Registrar o problema observado, a necessidade, a alternativa escolhida e como a mudança será testada.

## Validação física posterior

Quando houver equipamentos disponíveis:

1. Repetir conexão e reconexão do rádio.
2. Testar com dois ou mais nós e observar uplink/downlink.
3. Comparar identificadores e campos dos pacotes antes de classificar mensagens como duplicadas.
4. Avaliar o comportamento com PSKs e canais diferentes.
5. Atualizar o registro de equipamentos e os resultados observados.
6. Priorizar mudanças somente depois de identificar problemas concretos.

A ocorrência de três registros de uma mensagem em um ensaio anterior não determinou a causa. O hop limit do Meshtastic não significa que o broker deva receber uma cópia por hop; não implementar deduplicação com base somente nessa observação.

## Possibilidades futuras

### Autenticação MQTT e ACL

Necessárias antes de uma exposição pública mais ampla. A política deve limitar quem pode conectar, publicar e assinar tópicos. TLS sozinho não fornece autenticação MQTT de clientes nem autorização por tópico.

### Deduplicação

Avaliar somente depois de identificar publicações repetidas comparando IDs e campos relevantes. Definir identidade de pacote e janela de retenção com base no tráfego observado.

### Rate limiting

Avaliar limites por cliente ou outra unidade apropriada depois de conhecer o tráfego normal. Evitar valores arbitrários que bloqueiem tráfego legítimo.

### Observabilidade

Melhorar métricas e logs operacionais sem registrar conteúdo de mensagens desnecessariamente. Definir primeiro quais eventos são necessários para diagnóstico e segurança.

### Seleção de PSK por canal

A implementação atual utiliza uma única PSK configurada. Suporte a várias chaves exige decisão de arquitetura e testes com canais reais; não presumir que seja suportado.

### Persistência, moderação e bloqueio de clientes

São opções condicionais a uma necessidade operacional demonstrada. Não são requisitos do núcleo atual.

### Filtragem avançada de portnums e controle de hop

Não introduzir listas extensas de filtros nem reproduzir regras de encaminhamento de rádio no broker sem requisito técnico concreto e teste de compatibilidade.

## Critério de conclusão

O trabalho de software pode ser considerado concluído quando as correções aprovadas estão implementadas, os testes automatizados e de integração passam, vet e build passam, a imagem Docker é validada, a CI relevante está verde, a documentação reflete os limites reais e o diff foi revisado. A validação física multi-nó deve continuar marcada como pendente até que seja executada; não deve ser descrita como aprovada por inferência.

> Primeiro concluir o que pode ser demonstrado por código e testes reproduzíveis; depois validar o comportamento em uma malha física com múltiplos nós.
