## Descrição

Essa API tem como objeto a criação de um sistema de divulgação de costureiras.
A ideia é que as costureiras possam se cadastrar e divulgar seus locais de atendimento, bem como os serviços que oferecem.
Os clientes, por sua vez, poderão buscar por costureiras próximas a sua localização e entrar em contato com elas.

## Tecnologias

- Golang
- Docker
- PostgreSQL

## Funcionalidades

- Cadastro de costureiras
- [] Cadastro de clientes
- [] Busca de costureiras por localização
- [] Busca de costureiras por serviços oferecidos
- [] Avaliação de costureiras
- [] Comentários sobre costureiras
- [] Contato com costureiras

## Entidades

- Costureira

  - Nome
  - Contato
  - Serviços oferecidos
  - Local de atendimento
  - Avaliação
  - Comentários

- Cliente

  - Nome
  - Localização

## Regras de negócio

- Costureira

  - [] Não pode se cadastrar sem nome
  - [] Não pode se cadastrar sem contato
  - [] Não pode se cadastrar sem local de atendimento
  - [] Não pode se cadastrar sem serviços oferecidos

- Cliente
  - [] Não é obrigatório cadastrar para buscar costureiras
  - [] Pode buscar costureiras por localização e serviços oferecidos
  - [] Para buscar costureiras por localização, é necessário informar a localização e a distância máxima
  - [] Para avaliar costureiras, é necessário se cadastrar
  - [] Não pode se cadastrar sem nome
  - [] Não pode se cadastrar sem localização
  - [] Não pode avaliar costureiras sem comentários
