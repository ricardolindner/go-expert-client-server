# Projeto de Cotação de Dólar em Go

Este projeto consiste em uma aplicação client-server em Go.
* **Server**: Acessa uma API externa para obter a cotação do dólar, armazena os dados em um banco de dados SQLite e os retorna para o Client.
* **Client**: Realiza uma requisição HTTP para o servidor e salva a cotação recebida em um arquivo `cotacao.txt`.

## Requisitos
- Go 1.23.6 ou superior
- SQLite (necessário para a persistência das cotações)
- Dependência do `modernc.org/sqlite` para persistência no banco de dados SQLite.

## Como Rodar o Projeto
### 1. Clone o Repositório
Clone o repositório em sua máquina local:
```bash
git clone https://github.com/ricardolindner/go-expert-client-server.git
cd go-expert-client-server
```

### 2. Instale as Dependências
Dentro da pasta do projeto, instale as dependências do Go:
```bash
go mod tidy
```

### 3. Rodando o Server
O server irá acessar a API de cotação, salvar as informações no banco SQLite e expor as cotações através do endpoint /cotacao.
Para rodar o server, execute o seguinte comando:
```bash
cd server
go run server.go
```
O servidor estará rodando na porta `8080` e acessível via `http://localhost:8080/cotacao`.
* O timeout máximo para chamar a API de cotação do dólar é de 200ms;
* O timeout máximo para persistência dos dados no banco é de 10ms.

### 4. Rodando o Client
O client faz uma requisição HTTP ao server, obtém a cotação e a salva no arquivo cotacao.txt.
Para rodar o client, execute o seguinte comando em um novo terminal:
```bash
cd client
go run client.go
```
Após rodar o client, o arquivo cotacao.txt será gerado com a cotação do dólar.
* Timeout máximo para receber as informações do server é de 300ms.
