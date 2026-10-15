package pointer

import (
	"time"

	"github.com/google/uuid"
)

func CloneStringPointer(value *string) *string {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func CloneUUIDPointer(value *uuid.UUID) *uuid.UUID {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func CloneTimePointer(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func CloneUintPointer(value *uint) *uint {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}
