// Package querydsl parses and compiles the Elasticsearch-style Query DSL
// payload defined in especificacoes-projetos/backend/standard.md §1.4 into
// driver-agnostic SQL fragments. It knows nothing about any specific ORM or
// SQL driver: Compile returns a WHERE fragment, its bind args, an ORDER BY
// fragment and normalized LIMIT/OFFSET values, which the caller applies to
// whichever database layer it uses (GORM, database/sql, sqlx, ...).
package querydsl

// Clause is a single leaf query, e.g. {"term": {"name": "ACME"}}. The outer
// key is the operator (term, terms, match, range, exists); the inner map
// holds exactly one field name mapped to the operator's argument.
type Clause map[string]map[string]interface{}

// Bool groups clauses with boolean combinators, mirroring Elasticsearch's
// bool query: Must clauses are AND'ed, Should clauses are OR'ed together (as
// one AND'ed group), and MustNot clauses are AND NOT'ed.
type Bool struct {
	Must    []Clause `json:"must"`
	Should  []Clause `json:"should"`
	MustNot []Clause `json:"must_not"`
}

// Query wraps the Bool combinator under the "query" key of the payload.
type Query struct {
	Bool Bool `json:"bool"`
}

// SortField is a single sort criterion, e.g. {"created_at": "desc"}.
type SortField map[string]string

// Request is the full Query DSL payload accepted by POST /:entity/query and
// POST /:entity/trash/query across every API. The JSON contract (query,
// bool, must, should, must_not, sort, page, size) must stay stable.
type Request struct {
	Query Query       `json:"query"`
	Sort  []SortField `json:"sort"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}
