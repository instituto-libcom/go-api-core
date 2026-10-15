package pointer

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCloneStringPointer(t *testing.T) {
	if CloneStringPointer(nil) != nil {
		t.Fatal("expected nil")
	}
	original := "abc"
	cloned := CloneStringPointer(&original)
	if cloned == &original || *cloned != original {
		t.Fatal("expected distinct pointer with equal value")
	}
}

func TestCloneUUIDPointer(t *testing.T) {
	if CloneUUIDPointer(nil) != nil {
		t.Fatal("expected nil")
	}
	original := uuid.New()
	cloned := CloneUUIDPointer(&original)
	if cloned == &original || *cloned != original {
		t.Fatal("expected distinct pointer with equal value")
	}
}

func TestCloneTimePointer(t *testing.T) {
	if CloneTimePointer(nil) != nil {
		t.Fatal("expected nil")
	}
	original := time.Now()
	cloned := CloneTimePointer(&original)
	if cloned == &original || !cloned.Equal(original) {
		t.Fatal("expected distinct pointer with equal value")
	}
}

func TestCloneUintPointer(t *testing.T) {
	if CloneUintPointer(nil) != nil {
		t.Fatal("expected nil")
	}
	original := uint(7)
	cloned := CloneUintPointer(&original)
	if cloned == &original || *cloned != original {
		t.Fatal("expected distinct pointer with equal value")
	}
}
