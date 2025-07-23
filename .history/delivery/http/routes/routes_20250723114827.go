package routes

import "github.com/labstack/echo"

func (api *API) SetupRouter() {
	groupV1 := api.Echo.Group("/api/v1")
	tests := groupV1.Group("/tests")

	tests.GET("/", func(c echo.Context) error {

	})
	tests.POST("/upload", api.TestHandler.TestHandler)
}
