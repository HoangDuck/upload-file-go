package generates

import (
	"hrdept-web-service-2025/infrastructure/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func NewUUID() (string, error) {
	newId, err := uuid.NewUUID()
	if err != nil {
		logger.Error("error generate uuid", zap.Error(err))
		return "", err
	}
	return newId.String(), nil
}
