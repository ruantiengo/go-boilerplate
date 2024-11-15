# Boilerplate Go

Este é um boilerplate para projetos em Go, seguindo os princípios da Clean Architecture. Ele fornece uma estrutura básica para iniciar novos projetos de forma organizada e escalável.

## Dependencias 
- **sqlc**: Ferramenta para gerar código Go a partir de consultas SQL. [Instalação](https://docs.sqlc.dev/en/latest/overview/install.html)
- **migrate**: Ferramenta para gerenciar migrações de banco de dados. [Instalação](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- **postgres 17**: Banco de dados relacional utilizado no projeto. [Instalação](https://www.postgresql.org/download/)

## Estrutura de Pastas

A estrutura de pastas deste boilerplate é baseada na Clean Architecture, que promove a separação de responsabilidades e facilita a manutenção e evolução do código. Cada feature possui sua própria pasta, contendo suas entidades, casos de uso, interfaces, infraestrutura e requisitos.

```
/project-root
├── /cmd
│   └── /app            # Código da aplicação (onde o entrypoint principal reside)
├── /internal
│   ├── /transaction    # Feature "transaction"
│   │   ├── /domain    # Entidades e modelos do domínio
│   │   │   └── transaction.go  # Entidade Transaction
│   │   ├── /usecase   # Casos de uso relacionados à feature "transaction"
│   │   │   └── transaction_usecase.go
│   │   ├── /repository # Interfaces de repositório para a feature "transaction"
│   │   │   └── transaction_repository.go
│   │   ├── /service   # Serviços para a feature "transaction"
│   │   │   └── transaction_service.go
│   │   ├── /handler   # Manipuladores de requests (Controllers/Handlers)
│   │   │   └── transaction_handler.go
│   │   └── /test      # Testes para a feature "transaction"
│   │       ├── usecase_test.go
│   │       ├── repository_test.go
│   │       └── handler_test.go
│   ├── /database      # Configurações do banco de dados, migrações, etc.
│   │   └── migrations.go
│   ├── /config        # Configurações globais do sistema
│   │   └── config.go
│   └── /middleware    # Middlewares (autenticação, logging, etc.)
│       └── auth.go
├── /scripts           # Scripts auxiliares (migrações, seeders, etc.)
├── /docs              # Documentação do projeto
├── /web               # Framework web (pode ser gin, echo, etc.)
│   └── /router        # Roteamento HTTP
│       └── router.go
├── go.mod
├── go.sum
└── README.md
```

### Descrição das Pastas

- **cmd/**: Contém o ponto de entrada da aplicação. No exemplo, `app` inicializa e executa a aplicação.
- **internal/**: Contém o código interno da aplicação, dividido por features:
     - **transaction/**: Contém o código relacionado à feature de transações.
          - **domain/**: Contém as entidades de domínio da feature.
          - **usecase/**: Contém os casos de uso da feature, que representam as regras de negócio.
          - **repository/**: Contém as interfaces de repositório da feature.
          - **service/**: Contém os serviços da feature.
          - **handler/**: Contém os manipuladores de requests da feature.
          - **test/**: Contém os testes da feature.
     - **user/**: Contém o código relacionado à feature de usuários.
          - **domain/**: Contém as entidades de domínio da feature.
          - **usecase/**: Contém os casos de uso da feature, que representam as regras de negócio.
          - **repository/**: Contém as interfaces de repositório da feature.
          - **service/**: Contém os serviços da feature.
          - **handler/**: Contém os manipuladores de requests da feature.
          - **test/**: Contém os testes da feature.
     - **database/**: Contém as configurações do banco de dados, migrações, etc.
     - **config/**: Contém as configurações globais do sistema.
     - **middleware/**: Contém os middlewares (autenticação, logging, etc.).
- **scripts/**: Contém scripts auxiliares (migrações, seeders, etc.).
- **docs/**: Contém a documentação do projeto.
- **web/**: Contém o framework web (pode ser gin, echo, etc.) e o roteamento HTTP.

## Como Usar

1. Clone o repositório:
      ```sh
      git clone https://github.com/seu-usuario/boilerplate-go.git
      ```

2. Navegue até o diretório do projeto:
      ```sh
      cd boilerplate-go
      ```

3. Execute a aplicação:
      ```sh
      go run cmd/app/main.go
      ```

## Contribuição

Sinta-se à vontade para contribuir com melhorias e novas funcionalidades. Para isso, siga os passos abaixo:

1. Faça um fork do projeto.
2. Crie uma branch para sua feature:
      ```sh
      git checkout -b minha-feature
      ```
3. Commit suas mudanças:
      ```sh
      git commit -m 'Adiciona minha feature'
      ```
4. Envie para o repositório remoto:
      ```sh
      git push origin minha-feature
      ```
5. Abra um Pull Request.

## Licença

Este projeto está licenciado sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

