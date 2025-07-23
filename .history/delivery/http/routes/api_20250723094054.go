package routes

import "github.com/labstack/echo/v4"

type API struct {
	Echo        *echo.Echo
	TestHandler *tests.TestHandler
}
