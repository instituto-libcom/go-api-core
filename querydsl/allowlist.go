package querydsl

// FieldSpec describes how a filterable/sortable field exposed to the API
// maps onto the underlying table.
type FieldSpec struct {
	// Column is the real SQL column (or JSON column when JSONPath is true).
	// It is written verbatim into the generated SQL, so it must only ever
	// come from a hardcoded FieldSet — never from request input.
	Column string
	// JSONPath marks Column as a jsonb column that accepts dotted path
	// notation in the request, e.g. "json_contacts.email.pessoal".
	JSONPath bool
	// Sortable allows the field to be used in the request's "sort" list.
	Sortable bool
}

// FieldSet is the allowlist of fields an entity exposes to the Query DSL
// engine. It is the security boundary of the package: any field requested
// outside this set is rejected instead of reaching the database.
type FieldSet map[string]FieldSpec
