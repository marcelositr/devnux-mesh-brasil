# Uso de Inteligência Artificial no desenvolvimento

## Objetivo

Registrar como ferramentas de Inteligência Artificial podem participar do desenvolvimento do DevNux Mesh Brasil sem transferir para a ferramenta a responsabilidade pelo projeto.

Este documento trata do processo de desenvolvimento. As decisões de arquitetura, protocolo e operação continuam documentadas nos arquivos específicos:

- [Arquitetura](arquitetura.md);
- [Protocolo e criptografia](protocolo.md);
- [Configuração e execução](configuracao.md);
- [Guia de desenvolvimento](desenvolvimento.md).

## Papel da IA

Ferramentas de IA podem ser usadas para:

- analisar código;
- propor implementações;
- auxiliar na investigação de problemas;
- sugerir testes;
- revisar documentação;
- explicar APIs, protocolos e conceitos;
- executar tarefas repetitivas quando houver ferramentas apropriadas.

A IA é uma ferramenta de apoio. Ela não é autoridade sobre requisitos, arquitetura, segurança ou comportamento do projeto.

## Responsabilidade humana

Toda alteração relevante deve ser entendida e validada por uma pessoa responsável pelo projeto.

A aprovação humana não significa confiar cegamente no resultado gerado. Significa verificar:

- se a alteração atende ao requisito;
- se preserva decisões existentes;
- se não introduz comportamento não solicitado;
- se não cria dependências ou abstrações desnecessárias;
- se respeita licenças e direitos de terceiros;
- se os testes realmente cobrem o comportamento alterado.

## Regras para geração e alteração de código

A IA não deve:

- inventar requisitos;
- inventar resultados de testes;
- afirmar que executou uma ferramenta quando não executou;
- alterar o protocolo sem solicitação;
- substituir uma implementação funcional por outra apenas por preferência;
- introduzir arquitetura genérica sem necessidade;
- adicionar dependências sem justificar a necessidade;
- remover validações ou testes para fazer uma suíte passar;
- inserir credenciais, chaves privadas, tokens ou dados sensíveis no repositório.

Quando houver incerteza, a incerteza deve ser explicitada e a informação deve ser verificada antes de transformar a hipótese em código.

## Alterações incrementais

Prefira o seguinte ciclo:

1. entender o estado atual;
2. definir a mudança;
3. alterar uma responsabilidade por vez;
4. executar os testes;
5. revisar o diff;
6. registrar a mudança no Git;
7. somente então avançar.

Mudanças grandes devem ser divididas quando isso reduzir risco ou facilitar a revisão.

## Código de terceiros

Antes de incorporar código, documentação ou outro material de terceiros, verifique sua origem, licença e compatibilidade com a licença do projeto.

Não atribua autoria humana ou institucional que não possa ser comprovada.

## Autoria

O uso de IA durante o desenvolvimento não transforma a ferramenta em autora do projeto.

A autoria e a responsabilidade pelo software permanecem vinculadas às pessoas que definem, revisam, integram e mantêm o projeto.

O histórico Git deve continuar sendo a fonte de rastreabilidade das alterações realizadas.

## Verificação

Uma resposta ou implementação produzida com auxílio de IA não é considerada validada apenas por ter sido gerada.

Para código, a validação deve seguir [Guia de desenvolvimento](desenvolvimento.md), incluindo testes, análise estática, compilação e revisão do diff quando aplicável.

## Princípio

> IA pode acelerar o trabalho de desenvolvimento; não pode substituir a verificação do trabalho de desenvolvimento.
