# Evolução futura e hardening

## Objetivo

Registrar possibilidades de evolução sem transformar hipóteses em requisitos automáticos. A próxima etapa funcional é concluir os testes com hardware Meshtastic real, incluindo dois ou mais rádios. Consulte [Validação](validacao.md) e [Equipamentos e ensaios](equipamentos.md).

## Correções técnicas identificadas na revisão estática

Os itens abaixo são observações de leitura do código, não resultados de testes de execução. Devem ser avaliados como manutenção técnica separada da lista de funcionalidades futuras:

- **Erro ao iniciar o servidor:** o retorno de Serve é registrado numa goroutine, mas o processo principal continua aguardando SIGINT/SIGTERM. Se Serve terminar com erro, o processo pode permanecer ativo sem serviço funcional. Criar um teste para falha de bind e tratar o erro no ciclo de vida.
- **Logs de mensagens:** o hook registra o conteúdo integral das mensagens de texto em nível INFO. Avaliar redução ou remoção do conteúdo para limitar exposição de conversas e crescimento dos logs.
- **Teste de integração MQTT:** os testes do hook não substituem um cenário com dois clientes MQTT conectados, um publicando e outro recebendo a publicação. Cobrir esse fluxo em teste automatizado.
- **Limites de entrada:** verificar e documentar os limites efetivos de tamanho de publicação e conexões na biblioteca MQTT utilizada; não presumir que existam limites específicos configurados pelo projeto.

Esses itens devem ser confirmados por testes antes de declarar uma correção concluída. A revisão estática não executou o broker nem a suíte de testes.

## Regra de evolução

As funcionalidades listadas abaixo são possibilidades, não compromissos. Não implementar um item apenas porque aparece neste documento. Registrar o problema observado, a necessidade, a alternativa escolhida e como a mudança será testada.

## Sequência de validação física

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

## Critério para sair da fase de observação

Antes de adotar novas regras de tráfego, verificar estabilidade com hardware, comportamento de uplink/downlink, necessidade real de deduplicação, impacto na compatibilidade e forma de testar a mudança. Itens que não demonstrarem necessidade devem ser descartados.

> Primeiro compreender o comportamento real; depois implementar controles justificados por evidência.