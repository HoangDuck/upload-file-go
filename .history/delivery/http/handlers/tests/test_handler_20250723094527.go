package tests

import (
	"fmt"
	"net/http"
	"sound_qr_services/domain/usecases"

	"github.com/labstack/echo"
)

type TestHandler struct {
	UploadFileUseCase usecases.UploadFileUsecase
}

func (t *TestHandler) TestHandler(c echo.Context) error {
	//Check role
	_, err := t.PermissionUseCase.CheckAdminRole(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}
	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to get file from request",
		})
	}

	// Validate the file
	if err = uploadfiles.ValidateFile(file); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
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
