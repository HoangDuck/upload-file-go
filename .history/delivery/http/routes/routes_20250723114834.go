package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (api *API) SetupRouter() {
	groupV1 := api.Echo.Group("/api/v1")
	tests := groupV1.Group("/tests")

	tests.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello World")
	})
	tests.POST("/upload", api.TestHandler.TestHandler)
}
