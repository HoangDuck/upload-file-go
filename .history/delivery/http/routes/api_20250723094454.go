package routes

import (
	"sound_qr_services/delivery/http/handlers/tests"

	"github.com/labstack/echo/v4"
)

type API struct {
	Echo        *echo.Echo
	TestHandler *tests.TestHandler
}
