package tests

import (
	"fmt"
	"net/http"
	"sound_qr_services/domain/usecases"

	"github.com/labstack/echo/v4"
)

type TestHandler struct {
	UploadFileUseCase usecases.UploadFileUsecase
}

func (t *TestHandler) TestHandler(c echo.Context) error {
	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to get file from request",
		})
	}

	fileName, err := t.UploadFileUseCase.UploadFile(file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Error uploading file: %v", err),
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "File upload successfully",
		"name":    fileName,
	})
}
