# Segurança

## Relato responsável

Não publique vulnerabilidades de segurança como issues públicas.

Para um problema que possa expor credenciais, permitir acesso indevido, comprometer o broker, afetar criptografia ou causar execução não autorizada, utilize um canal privado de contato com o mantenedor antes da divulgação pública.

## Ao relatar

Inclua, quando possível:

- descrição do problema;
- impacto;
- versão ou commit afetado;
- passos mínimos para reprodução;
- evidências técnicas;
- mitigação conhecida.

Não inclua chaves privadas, PSKs reais, tokens ou outras credenciais.

## Escopo e limites conhecidos

O broker utiliza TLS no transporte e processamento criptográfico de pacotes Meshtastic. Alterações nessas áreas devem ser tratadas como mudanças sensíveis e acompanhadas de testes específicos.

A configuração atual não implementa autenticação MQTT nem autorização por tópico. O broker não deve ser exposto à Internet nem a redes não confiáveis antes de haver controles de acesso apropriados.

A negociação TLS, isoladamente, não garante a identidade do servidor: o cliente precisa validar o certificado. No teste inicial com hardware foi usado certificado autoassinado e o cliente estava configurado para não validar o certificado do servidor. Esse teste confirmou a conexão cifrada, mas não a validação da identidade do servidor.

O estado e os limites dos testes físicos estão descritos em [Validação](docs/validacao.md). Medidas adicionais de hardening devem seguir os critérios de [Evolução futura e hardening](docs/evolucao.md), após observação do comportamento real da rede.
