package utils

import (
	"fmt"

	"github.com/google/uuid"
)

func GenerateID(pre string) string {
	return fmt.Sprintf(pre, "_%s", uuid.New().String())
}

func GenerateUserID() string {
	return fmt.Sprintf("uid_%s", uuid.New().String())
}

func IsValidUserID(id string) bool {
	if len(id) < 5 {
		return false
	}
	return id[:4] == "uid_"
}
