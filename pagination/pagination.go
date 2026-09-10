package pagination

import (
	"fmt"
	"strconv"
	"strings"
)

// Params struct representing the parsed pagination parameters
type Params struct {
	Page int
	Size int
	Sort string
}

// Meta struct representing the pagination metadata for the JSON response
type Meta struct {
	Page          int `json:"page"`
	Size          int `json:"size"`
	TotalElements int `json:"total_elements"`
	TotalPages    int `json:"total_pages"`
}

// Parse parses pagination parameters from raw string inputs
func Parse(pageStr, sizeStr, sortStr string) (Params, error) {
	page := 0
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil {
			return Params{}, fmt.Errorf("parâmetro de página inválido: deve ser um inteiro")
		}
		if p < 0 {
			return Params{}, fmt.Errorf("parâmetro de página inválido: deve ser maior ou igual a 0")
		}
		page = p
	}

	size := 25
	if sizeStr != "" {
		s, err := strconv.Atoi(sizeStr)
		if err != nil {
			return Params{}, fmt.Errorf("parâmetro de tamanho inválido: deve ser um inteiro")
		}
		
		validSizes := map[int]bool{
			25:   true,
			50:   true,
			100:  true,
			250:  true,
			1000: true,
		}

		if !validSizes[s] {
			return Params{}, fmt.Errorf("parâmetro de tamanho inválido: valores permitidos são 25, 50, 100, 250, 1000")
		}
		size = s
	}

	sort := ""
	if sortStr != "" {
		parts := strings.Split(sortStr, ":")
		if len(parts) != 2 {
			return Params{}, fmt.Errorf("parâmetro de ordenação inválido: formato deve ser 'campo:asc' ou 'campo:desc'")
		}
		
		direction := strings.ToLower(parts[1])
		if direction != "asc" && direction != "desc" {
			return Params{}, fmt.Errorf("parâmetro de ordenação inválido: direção deve ser 'asc' ou 'desc'")
		}
		sort = sortStr
	}

	return Params{
		Page: page,
		Size: size,
		Sort: sort,
	}, nil
}
