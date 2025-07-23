package routes

import (
	"sound_qr_services/delivery/http/handlers/tests"
	"sound_qr_services/domain/usecases"

	"github.com/labstack/echo/v4"
)

type API struct {
	Echo        *echo.Echo
	TestHandler *tests.TestHandler
}

func NewAPI(echo *echo.Echo) {

	uploadFileUseCase := usecases.NewUploadFileUsecase()
	testHandler := &tests.TestHandler{
		UploadFileUseCase: uploadFileUseCase,
	}
	// Initialize the API struct with the provided Echo instance, SQL instance, and config
	api := &API{
		Echo:        echo,
		TestHandler: testHandler,
	}
	api.SetupRouter()
}
