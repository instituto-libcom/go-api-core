package querydsl

import (
	"fmt"
	"regexp"
	"strings"
)

var jsonSegmentPattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

var operators = map[string]func(field string, value interface{}, fields FieldSet) (string, []interface{}, error){
	"term":   buildTerm,
	"terms":  buildTerms,
	"match":  buildMatch,
	"range":  buildRange,
	"exists": buildExists,
}

// clauseCondition translates a single Clause — exactly one operator
// targeting exactly one field — into a parameterized SQL fragment.
func clauseCondition(clause Clause, fields FieldSet) (string, []interface{}, error) {
	if len(clause) != 1 {
		return "", nil, fmt.Errorf("cada cláusula deve declarar exatamente um operador")
	}

	for op, target := range clause {
		builder, ok := operators[op]
		if !ok {
			return "", nil, fmt.Errorf("operador não suportado: %q", op)
		}
		if len(target) != 1 {
			return "", nil, fmt.Errorf("operador %q deve ter exatamente um campo", op)
		}

		for field, value := range target {
			return builder(field, value, fields)
		}
	}

	return "", nil, fmt.Errorf("cláusula vazia")
}

func buildTerm(field string, value interface{}, fields FieldSet) (string, []interface{}, error) {
	column, err := resolveColumn(field, fields)
	if err != nil {
		return "", nil, err
	}
	return column + " = ?", []interface{}{value}, nil
}

func buildTerms(field string, value interface{}, fields FieldSet) (string, []interface{}, error) {
	column, err := resolveColumn(field, fields)
	if err != nil {
		return "", nil, err
	}

	values, ok := value.([]interface{})
	if !ok || len(values) == 0 {
		return "", nil, fmt.Errorf("terms requer uma lista não vazia para o campo %q", field)
	}
	return column + " IN ?", []interface{}{values}, nil
}

func buildMatch(field string, value interface{}, fields FieldSet) (string, []interface{}, error) {
	column, err := resolveColumn(field, fields)
	if err != nil {
		return "", nil, err
	}

	text, ok := value.(string)
	if !ok {
		return "", nil, fmt.Errorf("match requer um valor string para o campo %q", field)
	}
	return column + " ILIKE ?", []interface{}{"%" + text + "%"}, nil
}

var rangeOperators = []struct {
	key string
	sql string
}{
	{"gt", ">"}, {"gte", ">="}, {"lt", "<"}, {"lte", "<="},
}

func buildRange(field string, value interface{}, fields FieldSet) (string, []interface{}, error) {
	column, err := resolveColumn(field, fields)
	if err != nil {
		return "", nil, err
	}

	bounds, ok := value.(map[string]interface{})
	if !ok || len(bounds) == 0 {
		return "", nil, fmt.Errorf("range requer um objeto com gt/gte/lt/lte para o campo %q", field)
	}

	var fragments []string
	var args []interface{}
	for _, op := range rangeOperators {
		v, present := bounds[op.key]
		if !present {
			continue
		}
		fragments = append(fragments, column+" "+op.sql+" ?")
		args = append(args, v)
	}
	if len(fragments) == 0 {
		return "", nil, fmt.Errorf("range requer ao menos um de gt/gte/lt/lte para o campo %q", field)
	}

	return strings.Join(fragments, " AND "), args, nil
}

// buildExists ignores its value — {"exists": {"field": true}} only cares
// that "field" is the key, mirroring the field->argument shape of every
// other operator instead of Elasticsearch's {"exists": {"field": "name"}}.
func buildExists(field string, _ interface{}, fields FieldSet) (string, []interface{}, error) {
	column, err := resolveColumn(field, fields)
	if err != nil {
		return "", nil, err
	}
	return column + " IS NOT NULL", nil, nil
}

// resolveColumn translates a request field name into a trusted SQL column
// expression using fields as the allowlist. Plain fields must match a
// FieldSet entry exactly; jsonb fields are addressed as
// "<allowlisted-base>.<path...>", validating every path segment.
func resolveColumn(field string, fields FieldSet) (string, error) {
	if spec, ok := fields[field]; ok {
		if spec.JSONPath {
			return "", fmt.Errorf("campo %q requer um caminho jsonb, ex.: %q", field, field+".chave")
		}
		return spec.Column, nil
	}

	if idx := strings.Index(field, "."); idx > 0 {
		base, path := field[:idx], field[idx+1:]
		spec, ok := fields[base]
		if !ok || !spec.JSONPath {
			return "", fmt.Errorf("campo não permitido: %q", field)
		}

		segments := strings.Split(path, ".")
		for _, segment := range segments {
			if !jsonSegmentPattern.MatchString(segment) {
				return "", fmt.Errorf("caminho jsonb inválido em %q", field)
			}
		}

		return fmt.Sprintf("%s#>>'{%s}'", spec.Column, strings.Join(segments, ",")), nil
	}

	return "", fmt.Errorf("campo não permitido: %q", field)
}
