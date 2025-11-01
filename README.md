# Go Note Pad API

Uma API RESTful simples para anotações, construída com Go, seguindo a arquitetura MVC e utilizando um banco de dados MySQL.

## Configuração

### 1. Banco de Dados

Crie um banco de dados MySQL e um usuário com privilégios para acessá-lo. Em seguida, crie a tabela `notes` usando a seguinte instrução SQL:

```sql
CREATE TABLE notes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### 2. Variáveis de Ambiente

Defina as seguintes variáveis de ambiente para configurar a conexão com o banco de dados:

```bash
export DB_USER="seu_usuario_do_banco_de_dados"
export DB_PASSWORD="sua_senha_do_banco_de_dados"
export DB_HOST="seu_host_do_banco_de_dados"
export DB_PORT="sua_porta_do_banco_de_dados"
export DB_NAME="seu_nome_do_banco_de_dados"
```

### 3. Executar a Aplicação

```bash
go run main.go
```

A API estará disponível em `http://localhost:8080`.

## Endpoints da API

### Listar todas as notas

- **GET** `/notes`

**Exemplo com cURL:**
```bash
curl -X GET http://localhost:8080/notes
```

### Obter uma única nota

- **GET** `/notes/{id}`

**Exemplo com cURL:**
```bash
curl -X GET http://localhost:8080/notes/1
```

### Criar uma nova nota

- **POST** `/notes`

**Exemplo com cURL:**
```bash
curl -X POST http://localhost:8080/notes -H "Content-Type: application/json" -d '{"title": "Minha Primeira Nota", "content": "Este é o conteúdo."}'
```

### Atualizar uma nota existente

- **PUT** `/notes/{id}`

**Exemplo com cURL:**
```bash
curl -X PUT http://localhost:8080/notes/1 -H "Content-Type: application/json" -d '{"title": "Título Atualizado", "content": "Conteúdo atualizado."}'
```

### Excluir uma nota

- **DELETE** `/notes/{id}`

**Exemplo com cURL:**
```bash
curl -X DELETE http://localhost:8080/notes/1
```
