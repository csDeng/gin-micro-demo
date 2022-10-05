package utils

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

func GetUUId() string {
	return fmt.Sprintf("%s", uuid.NewV4())
}
