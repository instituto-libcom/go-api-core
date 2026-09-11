package softdelete

import (
	"time"

	"github.com/google/uuid"
)

// Fields contém os campos padrões para exclusão lógica e auditoria.
// Ele foi projetado para ser embutido (embedded) nos models.
type Fields struct {
	Deleted         bool       `json:"deleted" db:"deleted" gorm:"column:deleted;default:false;not null"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" db:"deleted_at" gorm:"column:deleted_at"`
	DeletedUserUUID *uuid.UUID `json:"deleted_user_uuid,omitempty" db:"deleted_user_uuid" gorm:"column:deleted_user_uuid;type:uuid"`
}
