# Configuração e execução

## Configuração padrão

O broker utiliza:

| Variável | Padrão | Finalidade |
|---|---|---|
| `MQTT_LISTEN_ADDRESS` | `:8883` | endereço do listener MQTT |
| `MQTT_CERTIFICATE_FILE` | `certificate.pem` | certificado TLS |
| `MQTT_PRIVATE_KEY_FILE` | `private.key` | chave privada TLS |
| `MESHTASTIC_DEFAULT_PSK` | PSK padrão do projeto | chave para descriptografia |

A PSK de ambiente deve ser uma string hexadecimal que resulte em 16 ou 32 bytes.

## Limites de entrada e conexões

O projeto usa Mochi MQTT v2.7.9 e não substitui os limites padrão da biblioteca por valores próprios. A documentação da versão indica:

- `MaximumPacketSize: 0`: não impõe limite máximo de pacote MQTT.
- `MaximumClients`: o valor padrão não constitui um teto operacional prático para esta implantação.
- `MaximumClientWritesPending: 8192`: capacidade padrão da fila de escritas pendentes de cada cliente; o Mochi alerta que valores altos aumentam o consumo de memória conforme cresce o número de clientes.

Referências: [API e opções do Mochi MQTT v2.7.9](https://pkg.go.dev/github.com/mochi-mqtt/server/v2@v2.7.9), seção *Server Options and Capabilities*. Não há variáveis de ambiente do projeto para ajustar esses limites.

Mantenha o listener em uma rede controlada e não considere a configuração atual endurecida para clientes não confiáveis. A meta inicial de planejamento é uma rede de até 200 rádios, mas isso não é uma capacidade comprovada nem significa 200 conexões MQTT: rádios podem publicar por meio de gateways. O limite de pacote MQTT também não deve ser confundido com o limite de payload LoRa; o broker recebe um `ServiceEnvelope` protobuf encapsulado em MQTT. Antes de definir valores explícitos, registre tamanhos reais dos envelopes e frequência de publicação em um ensaio representativo, estabeleça a topologia de gateways e meça memória, CPU e latência sob carga. O procedimento está detalhado em [Perfil de carga](perfil-de-carga.md). Não escolher limites arbitrários que possam rejeitar tráfego Meshtastic legítimo.

## TLS

O listener exige TLS 1.2 ou superior.

O certificado e a chave privada são carregados por:

- `MQTT_CERTIFICATE_FILE`;
- `MQTT_PRIVATE_KEY_FILE`.

A chave privada nunca deve ser adicionada ao repositório. O `.gitignore` e o `.dockerignore` excluem `certs/`, arquivos de ambiente locais e extensões comuns de certificados/chaves. O `.dockerignore` preserva `.env.example` como exemplo não secreto. Para verificar as exclusões do contexto Docker, execute `bash scripts/test-dockerignore.sh` com Docker instalado.

A CI testa o acesso à chave montada com permissão `0640`, usando o GID do grupo primário do usuário `devnux` no container. Em uma implantação real, reproduza o princípio de menor privilégio com o grupo configurado no host; o GID numérico precisa corresponder a um grupo efetivo do processo no container. A CI não substitui a verificação das permissões e da propriedade no host de implantação.

TLS cifra o transporte entre cliente e broker, mas a autenticação do servidor depende de o cliente validar o certificado apresentado. Um teste com certificado autoassinado e validação de certificado desabilitada confirma a negociação TLS, mas não comprova a identidade do servidor. Para ambientes além de testes locais, configure uma cadeia de confiança adequada e verifique que o cliente valida o certificado.

## Execução

O `go.mod` declara Go 1.27.0 como versão mínima e fixa a toolchain de validação em Go 1.27.2. O estágio de build do Docker também utiliza `golang:1.27.2-alpine`. Atualizações de toolchain devem ser acompanhadas da regressão completa.

Depois de obter as dependências:

```bash
go test ./...
go vet ./...
go test -race ./...
go run "golang.org/x/vuln/cmd/govulncheck@$(tr -d '[:space:]' < .github/govulncheck-version)" ./...
go build -o /tmp/devnux-mesh-broker ./cmd/broker
```

Para iniciar diretamente durante o desenvolvimento:

```bash
go run ./cmd/broker
```

O processo encerra de forma controlada ao receber SIGINT ou SIGTERM.

## Docker

O projeto possui um `Dockerfile` de múltiplas etapas que compila o broker e cria uma imagem de runtime mínima. As imagens-base do compilador Go e do runtime Alpine são fixadas por digest SHA-256, além das tags de versão. O Dependabot monitora atualizações mensais do ecossistema Docker; atualizações de digest devem passar pela CI completa.

### Imagem publicada

O GitHub Actions publica a imagem no GitHub Container Registry em:

```text
ghcr.io/marcelositr/devnux-mesh-brasil
```

A tag `latest` acompanha a branch `main`. Tags Git no formato `vX.Y.Z` também geram tags de versão no registry.

Para obter a imagem:

```bash
docker pull ghcr.io/marcelositr/devnux-mesh-brasil:latest
```

Construção local:

```bash
docker build -t devnux-mesh-brasil .
```

A imagem executa o broker como usuário sem privilégios. O certificado e a chave privada devem ser fornecidos em tempo de execução, fora da imagem.

**Permissões da chave privada em Linux:** não use permissões `0644` em uma chave real. Uma opção é restringir a chave a `root` e a um grupo dedicado (`0640`) e iniciar o container com o GID desse grupo como grupo suplementar (`--group-add`). O diretório montado também precisa permitir travessia ao grupo. Mantenha o bind mount somente leitura (`:ro`) e confirme a propriedade, o GID e as permissões no host. O usuário `devnux` continua sem privilégios; não torne a chave globalmente legível para contornar problemas de acesso.

Estrutura esperada no host:

```text
certs/
├── certificate.pem
└── private.key
```

Execução da imagem publicada, limitada ao próprio host:

```bash
docker run --rm \
  -p 127.0.0.1:8883:8883 \
  -v "$PWD/certs:/etc/devnux/certs:ro" \
  ghcr.io/marcelositr/devnux-mesh-brasil:latest
```

Dentro do container, os padrões são:

- `MQTT_LISTEN_ADDRESS=:8883`;
- `MQTT_CERTIFICATE_FILE=/etc/devnux/certs/certificate.pem`;
- `MQTT_PRIVATE_KEY_FILE=/etc/devnux/certs/private.key`.

Para utilizar outra PSK:

```bash
docker run --rm \
  -p 127.0.0.1:8883:8883 \
  -v "$PWD/certs:/etc/devnux/certs:ro" \
  -e MESHTASTIC_DEFAULT_PSK="d4f1bb3a20290759f0bcffabcf4e6901" \
  ghcr.io/marcelositr/devnux-mesh-brasil:latest
```

A chave privada, o certificado e outros segredos não devem ser copiados para dentro da imagem.

### Teste com um rádio na rede local

Se um rádio precisar conectar ao broker através da LAN, publique a porta somente no endereço privado específico do host, em vez de expô-la indiscriminadamente em todas as interfaces. Exemplo ilustrativo — substitua o endereço pelo IP privado real do host:

```bash
-p 192.0.2.10:8883:8883
```

O IP usado pelo rádio também precisa corresponder ao endereço coberto pelo certificado se o cliente validar o certificado. O exemplo de endereço acima é reservado para documentação e não deve ser usado como configuração real.

## Exemplo de configuração

```bash
export MQTT_LISTEN_ADDRESS=":8883"
export MQTT_CERTIFICATE_FILE="certificate.pem"
export MQTT_PRIVATE_KEY_FILE="private.key"
export MESHTASTIC_DEFAULT_PSK="d4f1bb3a20290759f0bcffabcf4e6901"

go run ./cmd/broker
```

Use certificados apropriados ao ambiente em que o broker será executado.

## Segurança operacional

A configuração atual aceita conexões e inscrições MQTT sem autenticação adicional e não implementa ACL. A porta não deve ser exposta à Internet nem a redes não confiáveis nessa condição.

Antes de uma implantação pública, os requisitos de autenticação, autorização, controle de abuso e observabilidade deverão ser definidos com base no comportamento observado em hardware Meshtastic físico e real.

Não habilite downlink no canal LongFast com a PSK padrão através deste broker enquanto não existir uma política zero-hop testada. O broker atual encaminha publicações sem reescrever `hop_limit`; use uma PSK privada em canal de teste para validar a ponte física sem injetar tráfego irrestrito na malha pública.

Não implemente essas camadas apenas com base em hipóteses. As possibilidades e os critérios de decisão estão em [Evolução futura e hardening](evolucao.md).

## Verificação

Após alterações de configuração ou infraestrutura, execute a bateria completa descrita em [Guia de desenvolvimento](desenvolvimento.md). Os resultados da validação funcional, Docker e hardware estão registrados em [Validação](validacao.md).
