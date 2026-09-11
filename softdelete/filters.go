package softdelete

// NotDeleted retorna a cláusula SQL ou abstração genérica para filtrar registros ativos.
// Exclui da query qualquer registro que tenha sido "soft deleted".
func NotDeleted() string {
	return "deleted = false OR deleted IS NULL"
}

// OnlyDeleted retorna a cláusula SQL ou abstração genérica para filtrar registros na lixeira.
// Mantém na query apenas os registros que foram "soft deleted".
func OnlyDeleted() string {
	return "deleted = true"
}
