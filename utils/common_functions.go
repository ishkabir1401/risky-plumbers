package utils

import (
	"github.com/google/uuid"
)

func GenerateUUid() string {
	u := uuid.New()
	return u.String()
}
