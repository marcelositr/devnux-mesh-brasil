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

## TLS

O listener exige TLS 1.2 ou superior.

O certificado e a chave privada são carregados por:

- `MQTT_CERTIFICATE_FILE`;
- `MQTT_PRIVATE_KEY_FILE`.

A chave privada nunca deve ser adicionada ao repositório.

## Execução

Depois de obter as dependências:

```bash
go test ./...
go vet ./...
go build -o /tmp/devnux-mesh-broker ./cmd/broker
```

Para iniciar diretamente durante o desenvolvimento:

```bash
go run ./cmd/broker
```

O processo encerra de forma controlada ao receber SIGINT ou SIGTERM.

## Docker

O projeto possui um `Dockerfile` de múltiplas etapas que compila o broker e cria uma imagem de runtime mínima.

Construção:

```bash
docker build -t devnux-mesh-brasil .
```

A imagem executa o broker como usuário sem privilégios. O certificado e a chave privada devem ser fornecidos em tempo de execução, fora da imagem.

Estrutura esperada no host:

```text
certs/
├── certificate.pem
└── private.key
```

Execução:

```bash
docker run --rm \
  -p 8883:8883 \
  -v "$PWD/certs:/etc/devnux/certs:ro" \
  devnux-mesh-brasil
```

Dentro do container, os padrões são:

- `MQTT_LISTEN_ADDRESS=:8883`;
- `MQTT_CERTIFICATE_FILE=/etc/devnux/certs/certificate.pem`;
- `MQTT_PRIVATE_KEY_FILE=/etc/devnux/certs/private.key`.

Para utilizar outra PSK:

```bash
docker run --rm \
  -p 8883:8883 \
  -v "$PWD/certs:/etc/devnux/certs:ro" \
  -e MESHTASTIC_DEFAULT_PSK="d4f1bb3a20290759f0bcffabcf4e6901" \
  devnux-mesh-brasil
```

A chave privada, o certificado e outros segredos não devem ser copiados para dentro da imagem.

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

A configuração atual aceita conexões e inscrições MQTT sem autenticação adicional e não implementa ACL.

Isso é uma característica do escopo atual, não uma indicação de que o broker deva ser exposto sem proteção em produção.

Antes de uma implantação pública, devem ser definidos requisitos de autenticação, autorização, controle de abuso e observabilidade.

## Verificação

Após alterações de configuração ou infraestrutura, execute a bateria completa descrita em [Guia de desenvolvimento](desenvolvimento.md).
