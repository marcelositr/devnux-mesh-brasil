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

## Escopo

O broker utiliza TLS e processamento criptográfico de pacotes Meshtastic. Alterações nessas áreas devem ser tratadas como mudanças sensíveis e acompanhadas de testes específicos.

A configuração atual não implementa autenticação MQTT nem autorização por tópico. Isso deve ser considerado em qualquer implantação exposta a redes não confiáveis.
