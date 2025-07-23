package tests

type TestHandler struct {
	UploadFileUseCase   usecases.UploadFileUsecase
	PermissionUseCase   usecases.PermissionUseCase
	NotificationUseCase usecases.NotificationUseCase
	TestUseCase         usecases.TestUseCase
}
