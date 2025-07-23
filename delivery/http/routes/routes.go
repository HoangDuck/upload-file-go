package routes

func (api *API) SetupRouter() {
	groupV1 := api.Echo.Group("/api/v1")
	tests := groupV1.Group("/tests")
	tests.POST("/upload", api.TestHandler.TestHandler)
}
