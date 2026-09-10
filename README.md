# go-api-core

- [Funcionalidades](https://github.com/instituto-libcom/especificacoes-projetos/blob/master/backend/standard.md) :wrench:

Pacote Go compartilhado com utilitários de resposta JSON, padronização de erros e paginação REST.
Usado para garantir a consistência de contratos entre os microserviços e APIs da organização.

## Instalação

```bash
go get github.com/instituto-libcom/go-api-core
```

## Exemplos de Uso

### 1. Resposta de Sucesso Única

```go
package main

import (
	"net/http"
	"github.com/instituto-libcom/go-api-core/response"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	user := User{ID: 1, Name: "João da Silva"}
	response.Success(w, http.StatusOK, user)
}
```

### 2. Resposta de Sucesso Paginada e Paginação

```go
package main

import (
	"net/http"
	"github.com/instituto-libcom/go-api-core/response"
	"github.com/instituto-libcom/go-api-core/pagination"
)

func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Obtendo parâmetros da Query String
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")
	sortStr := r.URL.Query().Get("sort")

	// Fazendo o parse e sanitização
	params, err := pagination.Parse(pageStr, sizeStr, sortStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	users := []User{{ID: 1, Name: "João"}, {ID: 2, Name: "Maria"}}
	
	meta := pagination.Meta{
		Page:          params.Page,
		Size:          params.Size,
		TotalElements: 150,
		TotalPages:    (150 / params.Size),
	}

	response.Paginated(w, http.StatusOK, users, meta)
}
```

### 3. Tratamento e Resposta de Erros

O pacote de erros (`apperror`) ajuda a padronizar e mapear códigos de status HTTP apropriados:

```go
package main

import (
	"net/http"
	"github.com/instituto-libcom/go-api-core/response"
	"github.com/instituto-libcom/go-api-core/apperror"
)

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Exemplo lançando um erro de Não Encontrado
	err := apperror.NewNotFound("USER_NOT_FOUND", "O usuário solicitado não existe")

	// Retornando a resposta padrão
	response.Error(w, err.StatusCode, err.Code, err.Message)
}
```
