package querydsl

import (
	"errors"
	"strings"
	"testing"
)

func testFields() FieldSet {
	return FieldSet{
		"name":          {Column: "name", Sortable: true},
		"status":        {Column: "status"},
		"created_at":    {Column: "created_at", Sortable: true},
		"rg_ie":         {Column: "rg_ie"},
		"json_contacts": {Column: "json_contacts", JSONPath: true},
	}
}

func mustCompile(t *testing.T, req Request) Compiled {
	t.Helper()

	compiled, err := Compile(req, testFields())
	if err != nil {
		t.Fatalf("Compile returned error: %v", err)
	}
	return compiled
}

func TestCompileTermFiltersByExactEquality(t *testing.T) {
	compiled := mustCompile(t, Request{Query: Query{Bool: Bool{Must: []Clause{
		{"term": {"name": "ACME"}},
	}}}})

	if compiled.Where != "(name = ?)" {
		t.Fatalf("unexpected WHERE: %q", compiled.Where)
	}
	if len(compiled.Args) != 1 || compiled.Args[0] != "ACME" {
		t.Fatalf("unexpected args: %+v", compiled.Args)
	}
}

func TestCompileMustAppliesAND(t *testing.T) {
	compiled := mustCompile(t, Request{Query: Query{Bool: Bool{Must: []Clause{
		{"term": {"name": "ACME"}},
		{"term": {"status": "ativo"}},
	}}}})

	if compiled.Where != "(name = ?) AND (status = ?)" {
		t.Fatalf("expected AND between must clauses, got %q", compiled.Where)
	}
	if len(compiled.Args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(compiled.Args))
	}
}

func TestCompileShouldAppliesGroupedOR(t *testing.T) {
	compiled := mustCompile(t, Request{Query: Query{Bool: Bool{Should: []Clause{
		{"term": {"name": "ACME"}},
		{"term": {"name": "Wayne"}},
	}}}})

	if compiled.Where != "((name = ?) OR (name = ?))" {
		t.Fatalf("expected grouped OR for should clauses, got %q", compiled.Where)
	}
}

func TestCompileMustNotExcludesMatches(t *testing.T) {
	compiled := mustCompile(t, Request{Query: Query{Bool: Bool{MustNot: []Clause{
		{"term": {"status": "bloqueado"}},
	}}}})

	if compiled.Where != "NOT (status = ?)" {
		t.Fatalf("expected NOT-wrapped must_not condition, got %q", compiled.Where)
	}
}

func TestCompileMustNotWrapsCompoundConditionCorrectly(t *testing.T) {
	compiled := mustCompile(t, Request{Query: Query{Bool: Bool{MustNot: []Clause{
		{"range": {"created_at": map[string]interface{}{"gte": "2026-01-01", "lte": "2026-12-31"}}},
	}}}})

	if compiled.Where != "NOT (created_at >= ? AND created_at <= ?)" {
		t.Fatalf("expected NOT to wrap the whole compound range condition, got %q", compiled.Where)
	}
}

func TestCompileCombinesMustShouldMustNotWithCorrectPrecedence(t *testing.T) {
	compiled := mustCompile(t, Request{Query: Query{Bool: Bool{
		Must:    []Clause{{"term": {"status": "ativo"}}},
		Should:  []Clause{{"term": {"name": "ACME"}}, {"term": {"name": "Wayne"}}},
		MustNot: []Clause{{"term": {"rg_ie": "123"}}},
	}}})

	expected := "(status = ?) AND ((name = ?) OR (name = ?)) AND NOT (rg_ie = ?)"
	if compiled.Where != expected {
		t.Fatalf("unexpected precedence:\n got: %q\nwant: %q", compiled.Where, expected)
	}
	if len(compiled.Args) != 4 {
		t.Fatalf("expected 4 args, got %d", len(compiled.Args))
	}
}

func TestCompileMatchUsesCaseInsensitivePartialSearch(t *testing.T) {
	compiled := mustCompile(t, Request{Query: Query{Bool: Bool{Must: []Clause{
		{"match": {"name": "acm"}},
	}}}})

	if compiled.Where != "(name ILIKE ?)" {
		t.Fatalf("expected ILIKE partial match, got %q", compiled.Where)
	}
	if compiled.Args[0] != "%acm%" {
		t.Fatalf("expected %%acm%%, got %v", compiled.Args[0])
	}
}

func TestCompileRangeAppliesInclusiveBounds(t *testing.T) {
	compiled := mustCompile(t, Request{Query: Query{Bool: Bool{Must: []Clause{
		{"range": {"created_at": map[string]interface{}{"gte": "2026-01-01", "lte": "2026-12-31"}}},
	}}}})

	if compiled.Where != "(created_at >= ? AND created_at <= ?)" {
		t.Fatalf("expected inclusive range bounds, got %q", compiled.Where)
	}
}

func TestCompileTermsTranslatesToIN(t *testing.T) {
	compiled := mustCompile(t, Request{Query: Query{Bool: Bool{Must: []Clause{
		{"terms": {"status": []interface{}{"ativo", "pendente"}}},
	}}}})

	if compiled.Where != "(status IN ?)" {
		t.Fatalf("expected IN clause, got %q", compiled.Where)
	}
}

func TestCompileExistsRequiresNonNullColumn(t *testing.T) {
	compiled := mustCompile(t, Request{Query: Query{Bool: Bool{Must: []Clause{
		{"exists": {"rg_ie": true}},
	}}}})

	if compiled.Where != "(rg_ie IS NOT NULL)" {
		t.Fatalf("expected IS NOT NULL condition, got %q", compiled.Where)
	}
}

func TestCompileRejectsFieldOutsideAllowlist(t *testing.T) {
	_, err := Compile(Request{Query: Query{Bool: Bool{Must: []Clause{
		{"term": {"password_hash": "x"}},
	}}}}, testFields())
	if err == nil {
		t.Fatal("expected an error for a field outside the allowlist")
	}
}

func TestCompileTreatsValuesAsBindArgsNeverInterpolated(t *testing.T) {
	malicious := "'; DROP TABLE client; --"
	compiled := mustCompile(t, Request{Query: Query{Bool: Bool{Must: []Clause{
		{"term": {"name": malicious}},
	}}}})

	if strings.Contains(compiled.Where, "DROP TABLE") {
		t.Fatalf("expected the payload to never be interpolated into WHERE, got %q", compiled.Where)
	}
	if compiled.Where != "(name = ?)" {
		t.Fatalf("expected a plain bind placeholder, got %q", compiled.Where)
	}
	if len(compiled.Args) != 1 || compiled.Args[0] != malicious {
		t.Fatalf("expected the malicious value to travel only as a bind arg, got %+v", compiled.Args)
	}
}

func TestCompileSortValidatesDirectionAndAllowlist(t *testing.T) {
	compiled := mustCompile(t, Request{Sort: []SortField{{"name": "asc"}}})
	if compiled.OrderBy != "name ASC" {
		t.Fatalf("expected ORDER BY name ASC, got %q", compiled.OrderBy)
	}

	if _, err := Compile(Request{Sort: []SortField{{"name": "ascending"}}}, testFields()); err == nil {
		t.Fatal("expected an error for an invalid sort direction")
	}
	if _, err := Compile(Request{Sort: []SortField{{"password_hash": "asc"}}}, testFields()); err == nil {
		t.Fatal("expected an error for a sort field outside the allowlist")
	}
	if _, err := Compile(Request{Sort: []SortField{{"status": "asc"}}}, testFields()); err == nil {
		t.Fatal("expected an error for a sort field that is not marked Sortable")
	}
}

func TestCompileEmptyPayloadYieldsDefaultPagination(t *testing.T) {
	compiled := mustCompile(t, Request{})

	if compiled.Where != "" {
		t.Fatalf("expected no WHERE fragment for an empty payload, got %q", compiled.Where)
	}
	if compiled.Page != 0 || compiled.Size != 25 || compiled.Limit != 25 || compiled.Offset != 0 {
		t.Fatalf("unexpected default pagination: %+v", compiled)
	}
}

func TestCompilePaginationComputesLimitAndOffset(t *testing.T) {
	compiled := mustCompile(t, Request{Page: 2, Size: 50})

	if compiled.Limit != 50 || compiled.Offset != 100 {
		t.Fatalf("expected limit=50 offset=100, got limit=%d offset=%d", compiled.Limit, compiled.Offset)
	}
}

func TestCompileRejectsInvalidSize(t *testing.T) {
	if _, err := Compile(Request{Size: 99}, testFields()); err == nil {
		t.Fatal("expected an error for an invalid size")
	}
}

func TestCompileRejectsNegativePage(t *testing.T) {
	if _, err := Compile(Request{Page: -1}, testFields()); err == nil {
		t.Fatal("expected an error for a negative page")
	}
}

func TestCompileJSONPathFiltersNestedField(t *testing.T) {
	compiled := mustCompile(t, Request{Query: Query{Bool: Bool{Must: []Clause{
		{"term": {"json_contacts.email.pessoal": "a@b.com"}},
	}}}})

	if compiled.Where != "(json_contacts#>>'{email,pessoal}' = ?)" {
		t.Fatalf("expected jsonb path expression, got %q", compiled.Where)
	}
}

func TestCompileRejectsJSONPathWithInvalidSegment(t *testing.T) {
	_, err := Compile(Request{Query: Query{Bool: Bool{Must: []Clause{
		{"term": {"json_contacts.email; DROP TABLE client;--": "x"}},
	}}}}, testFields())
	if err == nil {
		t.Fatal("expected an error for an invalid jsonb path segment")
	}
}

func TestCompileRejectsBareJSONPathBaseField(t *testing.T) {
	_, err := Compile(Request{Query: Query{Bool: Bool{Must: []Clause{
		{"term": {"json_contacts": "a@b.com"}},
	}}}}, testFields())
	if err == nil {
		t.Fatal("expected an error when a jsonb field is referenced without a path")
	}
}

func TestCompileRejectsTooManyClauses(t *testing.T) {
	must := make([]Clause, MaxClauses+1)
	for i := range must {
		must[i] = Clause{"term": {"name": "ACME"}}
	}

	_, err := Compile(Request{Query: Query{Bool: Bool{Must: must}}}, testFields())
	if err == nil {
		t.Fatal("expected an error when clause count exceeds MaxClauses")
	}
}

func TestCompileRejectsUnsupportedOperator(t *testing.T) {
	_, err := Compile(Request{Query: Query{Bool: Bool{Must: []Clause{
		{"wildcard": {"name": "AC*"}},
	}}}}, testFields())
	if err == nil {
		t.Fatal("expected an error for an unsupported operator")
	}
}

// CT-02 (issue #67): should/must_not may be omitted entirely without error.
func TestCompileAllowsOmittedShouldAndMustNot(t *testing.T) {
	compiled := mustCompile(t, Request{Query: Query{Bool: Bool{
		Must: []Clause{{"term": {"name": "ACME"}}},
	}}})

	if compiled.Where != "(name = ?)" {
		t.Fatalf("unexpected WHERE with should/must_not omitted: %q", compiled.Where)
	}
}

func TestCompileReturnsPlainError(t *testing.T) {
	_, err := Compile(Request{Page: -1}, testFields())
	if err == nil {
		t.Fatal("expected an error")
	}
	// The package must not force any particular error type on callers —
	// consistent with pagination.Parse, it returns a plain error.
	var target interface{ StatusCode() int }
	if errors.As(err, &target) {
		t.Fatalf("expected a plain error, not one carrying an HTTP status code: %T", err)
	}
}
