# go-api-core

- [Funcionalidades](https://github.com/instituto-libcom/especificacoes-projetos/blob/master/backend/standard.md) :wrench:

Pacote Go compartilhado com utilitários de resposta JSON, padronização de erros e paginação REST.
Usado para garantir a consistência de contratos entre os microserviços e APIs da organização.

## Instalação

```bash
go get github.com/instituto-libcom/go-api-core@latest
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

### 4. Gerenciamento de Lixeira (Soft Delete)

O subpacote `softdelete` abstrai as regras de exclusão lógica e física, fornecendo campos padronizados e utilitários agnósticos que se integram facilmente a ORMs como GORM.

**Incorporando os campos padrão (Entity/Model):**

```go
package model

import "github.com/instituto-libcom/go-api-core/softdelete"

type Product struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `json:"name"`
	Price float64
	
	// Adiciona as colunas deleted, deleted_at e deleted_user_uuid
	softdelete.Fields 
}
```

**Exemplos de uso (com GORM):**

```go
package repository

import (
	"github.com/google/uuid"
	"github.com/instituto-libcom/go-api-core/softdelete"
	"gorm.io/gorm"
)

// Listar apenas ativos
func ListActive(db *gorm.DB) ([]Product, error) {
	var products []Product
	err := db.Where(softdelete.NotDeleted()).Find(&products).Error
	return products, err
}

// Listar lixeira
func ListTrashed(db *gorm.DB) ([]Product, error) {
	var products []Product
	err := db.Where(softdelete.OnlyDeleted()).Find(&products).Error
	return products, err
}

// Executar Soft Delete
func SoftDelete(db *gorm.DB, productIDs []uint, userUUID uuid.UUID) error {
	updates := softdelete.SoftDeleteData(userUUID)
	return db.Model(&Product{}).Where("id IN ?", productIDs).Updates(updates).Error
}

// Restaurar da Lixeira
func Restore(db *gorm.DB, productIDs []uint) error {
	updates := softdelete.RestoreData()
	return db.Model(&Product{}).Where("id IN ?", productIDs).Updates(updates).Error
}
```
