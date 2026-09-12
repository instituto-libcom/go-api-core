package querydsl

import (
	"fmt"
	"strings"
)

// MaxClauses bounds the total number of must/should/must_not clauses accepted
// in a single request, to keep pathological payloads from generating
// arbitrarily large queries.
const MaxClauses = 50

// allowedSizes fixes the page-size options every API must honor per
// especificacoes-projetos/backend/standard.md §1.3.
var allowedSizes = map[int]bool{25: true, 50: true, 100: true, 250: true, 1000: true}

// Compiled is the driver-agnostic result of compiling a Request: a WHERE
// fragment with its positional ("?") bind args, an ORDER BY fragment, and
// normalized pagination. Where and OrderBy are empty strings when the
// request declares no filters/sort, so callers should only append them to a
// query when non-empty.
type Compiled struct {
	Where   string
	Args    []interface{}
	OrderBy string
	Page    int
	Size    int
	Limit   int
	Offset  int
}

// Compile translates req into Compiled using fields as the allowlist of
// filterable/sortable columns. It never touches a database or ORM — the
// caller applies Where/Args, OrderBy, Limit and Offset to whatever query
// builder it uses (see the README for a GORM example).
func Compile(req Request, fields FieldSet) (Compiled, error) {
	b := req.Query.Bool
	if len(b.Must)+len(b.Should)+len(b.MustNot) > MaxClauses {
		return Compiled{}, fmt.Errorf("número de cláusulas excede o limite de %d", MaxClauses)
	}

	where, args, err := compileBool(b, fields)
	if err != nil {
		return Compiled{}, err
	}

	orderBy, err := compileSort(req.Sort, fields)
	if err != nil {
		return Compiled{}, err
	}

	page, size, err := normalizePage(req.Page, req.Size)
	if err != nil {
		return Compiled{}, err
	}

	return Compiled{
		Where:   where,
		Args:    args,
		OrderBy: orderBy,
		Page:    page,
		Size:    size,
		Limit:   size,
		Offset:  page * size,
	}, nil
}

// compileBool AND's must clauses, OR's should clauses together as a single
// group, and AND NOT's must_not clauses — each clause fragment is wrapped in
// parens so the composition is correct regardless of whether the fragment
// itself is a single comparison or a compound one (e.g. a range).
func compileBool(b Bool, fields FieldSet) (string, []interface{}, error) {
	var parts []string
	var args []interface{}

	for _, clause := range b.Must {
		condition, clauseArgs, err := clauseCondition(clause, fields)
		if err != nil {
			return "", nil, err
		}
		parts = append(parts, "("+condition+")")
		args = append(args, clauseArgs...)
	}

	if len(b.Should) > 0 {
		var orParts []string
		for _, clause := range b.Should {
			condition, clauseArgs, err := clauseCondition(clause, fields)
			if err != nil {
				return "", nil, err
			}
			orParts = append(orParts, "("+condition+")")
			args = append(args, clauseArgs...)
		}
		parts = append(parts, "("+strings.Join(orParts, " OR ")+")")
	}

	for _, clause := range b.MustNot {
		condition, clauseArgs, err := clauseCondition(clause, fields)
		if err != nil {
			return "", nil, err
		}
		parts = append(parts, "NOT ("+condition+")")
		args = append(args, clauseArgs...)
	}

	return strings.Join(parts, " AND "), args, nil
}

func compileSort(sort []SortField, fields FieldSet) (string, error) {
	if len(sort) == 0 {
		return "", nil
	}

	var parts []string
	for _, entry := range sort {
		if len(entry) != 1 {
			return "", fmt.Errorf("cada critério de sort deve conter exatamente um campo")
		}

		for field, direction := range entry {
			spec, ok := fields[field]
			if !ok || !spec.Sortable {
				return "", fmt.Errorf("campo de ordenação não permitido: %q", field)
			}

			normalized := strings.ToLower(strings.TrimSpace(direction))
			if normalized != "asc" && normalized != "desc" {
				return "", fmt.Errorf("direção de ordenação inválida para %q: use \"asc\" ou \"desc\"", field)
			}

			column, err := resolveColumn(field, fields)
			if err != nil {
				return "", err
			}
			parts = append(parts, column+" "+strings.ToUpper(normalized))
		}
	}

	return strings.Join(parts, ", "), nil
}

// normalizePage validates page/size and fills in the page=0/size=25
// defaults for an omitted (zero-valued) request field.
func normalizePage(page, size int) (int, int, error) {
	if page < 0 {
		return 0, 0, fmt.Errorf("page deve ser maior ou igual a 0")
	}
	if size == 0 {
		size = 25
	} else if !allowedSizes[size] {
		return 0, 0, fmt.Errorf("size inválido: valores permitidos são 25, 50, 100, 250, 1000")
	}
	return page, size, nil
}
