version: "2"
packages:
  - name: "queries"          # Nome do pacote Go gerado, que terá os repositórios e queries
    path: "db/queries"        # Diretório onde os arquivos SQL com queries estão localizados
    engine: "postgresql"      # Defina o tipo de banco de dados (ex: postgresql, mysql, sqlite)
    schema: "db/schema.sql"   # Arquivo SQL com o esquema do banco (opcional se o schema estiver no próprio banco)
    queries: "db/queries"     # Diretório com os arquivos SQL de queries
    gen:
      go:                     # Gera código Go para os modelos e queries
        out: "db/queries"      # Pasta onde os arquivos Go gerados serão armazenados
        package: "queries"     # Nome do pacote Go gerado para queries
      sql:                     # Gera o SQL para as tabelas (se necessário)
        out: "db/queries"      # Pasta onde os arquivos SQL gerados serão armazenados
    options:
      sql_package: "pgx/v4"    # Pacote SQL que será usado para conexão ao banco
