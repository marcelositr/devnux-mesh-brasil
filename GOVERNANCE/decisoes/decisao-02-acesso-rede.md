# Decisão 02: manter o broker em rede controlada

**Estado:** orientação operacional vigente enquanto não houver controles de acesso adequados.

A configuração atual não implementa autenticação MQTT nem ACL por tópico. TLS protege o transporte, mas não substitui autenticação do cliente nem garante que o cliente valide a identidade do servidor.

Até que os controles de acesso sejam projetados, implementados e testados, manter o listener em ambiente local ou rede controlada e não expô-lo à Internet ou a redes não confiáveis.

Esta orientação não aprova uma solução específica de autenticação. As alternativas devem ser avaliadas separadamente, com critérios de teste e documentação atualizada.

Referências: [Segurança](../../SECURITY.md) e [Configuração e execução](../../docs/configuracao.md).