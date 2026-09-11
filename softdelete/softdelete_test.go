package softdelete

import (
	"testing"

	"github.com/google/uuid"
)

func TestSoftDeleteData(t *testing.T) {
	userUUID := uuid.New()
	data := SoftDeleteData(userUUID)

	if deleted, ok := data["deleted"].(bool); !ok || !deleted {
		t.Errorf("Expected deleted to be true, got %v", data["deleted"])
	}

	if data["deleted_at"] == nil {
		t.Error("Expected deleted_at to not be nil")
	}

	if uuidVal, ok := data["deleted_user_uuid"].(uuid.UUID); !ok || uuidVal != userUUID {
		t.Errorf("Expected deleted_user_uuid to be %v, got %v", userUUID, data["deleted_user_uuid"])
	}
}

func TestRestoreData(t *testing.T) {
	data := RestoreData()

	if deleted, ok := data["deleted"].(bool); !ok || deleted {
		t.Errorf("Expected deleted to be false, got %v", data["deleted"])
	}

	if data["deleted_at"] != nil {
		t.Errorf("Expected deleted_at to be nil, got %v", data["deleted_at"])
	}

	if data["deleted_user_uuid"] != nil {
		t.Errorf("Expected deleted_user_uuid to be nil, got %v", data["deleted_user_uuid"])
	}
}

func TestConditions(t *testing.T) {
	active := NotDeleted()
	expectedActive := "deleted = false OR deleted IS NULL"
	if active != expectedActive {
		t.Errorf("Expected NotDeleted to be %q, got %q", expectedActive, active)
	}

	trashed := OnlyDeleted()
	expectedTrashed := "deleted = true"
	if trashed != expectedTrashed {
		t.Errorf("Expected OnlyDeleted to be %q, got %q", expectedTrashed, trashed)
	}
}
