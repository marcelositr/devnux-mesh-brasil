# Perfil de carga — planejamento e validação

## Objetivo e limite da evidência

A meta inicial de planejamento é uma rede de até **200 rádios**. Isso não é uma declaração de capacidade suportada nem implica 200 conexões MQTT simultâneas. Um rádio pode publicar diretamente no broker ou o tráfego pode chegar por meio de um gateway; os dois modelos têm perfis de conexão diferentes.

Este documento define como obter evidências antes de configurar limites. Não fixa taxas de publicação, limites de pacotes ou metas de latência sem medições do cenário real.

## Dados que devem ser coletados

Antes de uma execução de carga representativa, registrar:

- Topologia: rádios que se conectam diretamente, quantidade de gateways e número de clientes MQTT simultâneos.
- Versão do firmware Meshtastic e configuração MQTT de cada tipo de nó.
- Canais com uplink/downlink ativos, criptografia e opções de relatório habilitadas.
- Tamanho em bytes dos `ServiceEnvelope` protobuf recebidos, incluindo mínimo, mediana, percentil 95 e máximo observados.
- Publicações por minuto em operação normal e durante o maior pico observado.
- Tipos de tráfego presentes e distribuição por tipo de pacote.
- Recursos do host/container: CPU, memória disponível, limites do container e versão da imagem.
- Condições de rede e comportamento de reconexão dos clientes.

Não registrar conteúdo sensível das mensagens. Identificadores de teste e tamanhos de payload bastam para correlacionar publicações e medir entrega.

## Cenários de teste

Executar os cenários separadamente para identificar a causa de uma falha:

1. **Conexões ociosas:** simular até 200 clientes MQTT simultâneos como cenário de pior caso para conexão direta. Não interpretar esse cenário como equivalente a 200 rádios em uma topologia com gateways.
2. **Topologia com gateways:** repetir com o número de conexões correspondente à implantação planejada, encaminhando o volume agregado de tráfego medido.
3. **Tráfego normal:** reproduzir envelopes válidos capturados/anonimizados ou gerados com o mesmo tamanho e distribuição observados no ensaio físico.
4. **Pico observado:** reproduzir a maior taxa medida, sem multiplicadores arbitrários antes de conhecer o perfil real.
5. **Assinante lento:** verificar o comportamento quando um cliente recebe mais devagar do que o broker publica.
6. **Reconexões:** testar desconexões e reconexões simultâneas, observando estabilidade e recuperação.
7. **Entradas inválidas:** verificar descarte de protobuf inválido, envelopes incompletos e publicações MQTT que excedam o limite que venha a ser definido.

A execução precisa cobrir os ciclos periódicos de tráfego habilitados no firmware. Para a validação de release, realizar também um ensaio prolongado e registrar sua duração real.

## Métricas e resultados

Registrar para cada cenário:

- Número de clientes conectados e reconexões.
- Publicações enviadas, aceitas, ignoradas e efetivamente recebidas pelos assinantes.
- Latência ponta a ponta entre publicação e recebimento, quando o cliente de teste permitir correlacionar as mensagens.
- CPU, memória residente e pico de memória do processo/container.
- Erros de protocolo, timeouts, encerramentos inesperados e mensagens perdidas.
- Estado do container e logs relevantes, sem conteúdo privado das mensagens.

Guardar a configuração, o commit testado, a versão da imagem, o conjunto de dados de teste e os resultados. Não declarar capacidade com base apenas no número de conexões ou em um único ensaio.

## Critérios para decidir limites

Os valores de `MaximumPacketSize`, `MaximumClients` e `MaximumClientWritesPending` devem ser decididos somente depois de analisar as medições e os requisitos da implantação. Confirmar também o impacto das filas de escrita e dos buffers por cliente, descritos na documentação do Mochi MQTT v2.7.9 em [Server Options and Capabilities](https://pkg.go.dev/github.com/mochi-mqtt/server/v2@v2.7.9).

Cada limite candidato deve ter testes de fronteira (valor permitido e valor imediatamente acima), teste de regressão de tráfego Meshtastic legítimo e uma execução de carga. A capacidade declarada deve corresponder à topologia e aos recursos do host efetivamente testados.

## Estado atual

- Meta de planejamento: até 200 rádios.
- Topologia final: ainda não definida.
- Perfil de tráfego real com tamanhos e taxas medidos: pendente.
- Testes de carga e capacidade comprovada: pendentes.
- Validação física multi-rádio: pendente.
