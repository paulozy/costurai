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
- [x] Cadastro de costureiras
- [x] Login de costureiras
- [X] Busca de costureiras por localização
- [X] Busca de costureiras por serviços oferecidos
- [X] Avaliação de costureiras
- [X] Comentários sobre costureiras
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

  - [x] Não pode se cadastrar sem nome
  - [x] Não pode se cadastrar sem contato
  - [x] Não pode se cadastrar sem local de atendimento
  - [x] Não pode se cadastrar sem serviços oferecidos

- Cliente
  - [X] Não é obrigatório cadastrar para buscar costureiras
  - [X] Pode buscar costureiras por localização e serviços oferecidos
  - [X] Para buscar costureiras por localização, é necessário informar a localização e a distância máxima
  - [X] Para avaliar costureiras, é necessário se cadastrar
  - [X] Não pode se cadastrar sem nome
  - [X] Não pode se cadastrar sem localização
  - [X] Não pode avaliar costureiras sem comentários
