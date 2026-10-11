# Segurança

## Relato responsável

Não publique detalhes de vulnerabilidades como issues públicas. Use o recurso privado de reporte de vulnerabilidades do GitHub caso esteja habilitado para este repositório; caso contrário, contate o mantenedor por um canal privado já estabelecido antes de divulgar os detalhes. A configuração atual do repositório não documenta um endereço dedicado para relatos de segurança, portanto esse canal operacional ainda precisa ser definido antes de uma distribuição pública mais ampla.

Ao relatar, inclua, quando possível:

- descrição do problema e impacto;
- versão ou commit afetado;
- passos mínimos para reprodução;
- evidências técnicas e mitigação conhecida.

Não inclua chaves privadas, PSKs reais, tokens ou outras credenciais.

## Estado de segurança conhecido

O broker utiliza TLS no transporte e processamento criptográfico de pacotes Meshtastic. Alterações nessas áreas são sensíveis e devem incluir testes específicos.

A configuração atual não implementa autenticação MQTT nem autorização por tópico. Não exponha o broker à Internet ou a redes não confiáveis antes de implementar e validar controles de acesso apropriados.

A negociação TLS isoladamente não garante a identidade do servidor: o cliente precisa validar o certificado. No teste físico documentado foi usado certificado autoassinado e o cliente não validava o certificado do servidor. Esse ensaio confirmou a conexão cifrada, mas não a validação da identidade do servidor.

## Evidências e próximos passos

Os resultados físicos e seus limites estão em [Validação](docs/validacao.md); a matriz de hardware está em [Equipamentos e ensaios](docs/equipamentos.md). Medidas futuras de hardening devem seguir [Evolução futura e hardening](docs/evolucao.md) e ser justificadas por requisitos e evidências.