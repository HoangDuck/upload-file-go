package generates

import (
	"github.com/google/uuid"
)

func NewUUID() (string, error) {
	newId, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return newId.String(), nil
}
