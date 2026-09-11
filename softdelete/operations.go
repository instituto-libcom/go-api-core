package softdelete

import (
	"time"

	"github.com/google/uuid"
)

// SoftDeleteData retorna os campos e valores necessários para realizar um Soft Delete lógico.
// Ele define deleted = true, deleted_at = time.Now() e atribui o UUID do usuário que realizou a ação.
func SoftDeleteData(userUUID uuid.UUID) map[string]interface{} {
	return map[string]interface{}{
		"deleted":           true,
		"deleted_at":        time.Now(),
		"deleted_user_uuid": userUUID,
	}
}

// RestoreData retorna os campos e valores necessários para restaurar um registro da lixeira.
// Ele limpa os campos de deleção, tornando o registro ativo novamente.
func RestoreData() map[string]interface{} {
	return map[string]interface{}{
		"deleted":           false,
		"deleted_at":        nil,
		"deleted_user_uuid": nil,
	}
}
