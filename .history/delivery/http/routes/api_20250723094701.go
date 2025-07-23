package routes

import (
	"sound_qr_services/delivery/http/handlers/tests"

	"github.com/labstack/echo/v4"
)

type API struct {
	Echo        *echo.Echo
	TestHandler *tests.TestHandler
}

func NewAPI(echo *echo.Echo) {

	uploadFileUseCase := usecases.NewUploadFileUsecase()
	// Initialize the Email instance
	gmailService := gmail.GetEmailGoogleInstance(
		config.Email.SenderEmail,
		config.Email.SenderPass,
	)
	configurationUseCase := usecases.NewConfigurationUseCase(
		configuration_repository.NewConfigurationRepoImpl(sql),
	)
	//Initialize Company repositories and use cases

	thumbnailRepo := thumbnail_repository.NewThumbnailRepository(
		sql,
	)

	thumbnailUseCase := usecases.NewThumbnailUseCase(
		thumbnailRepo,
	)

	cronTaskRepo := crons_repository.NewCronsTaskRepoImpl(sql)
	cronTaskUseCase := usecases.NewCronTaskUseCase(cronTaskRepo)

	matchingJobCVRepo := matching_job_repository.NewMatchingJobRepository(sql)
	userRepo := user_repository.NewUserRepoImpl(sql)
	companyRepo := company_repository.NewCompanyRepoImpl(sql)
	companyUseCase := usecases.NewCompanyUseCase(companyRepo, userRepo, uploadFileUseCase)
	// Initialize user repositories and use cases
	authVerifyRepo := auth_verify_repository.NewAuthVerificationRepositoryImpl(sql)
	funcRepo := functions_repository.NewFunctionsRepoImpl(sql)
	roleRepo := role_repository.NewRoleRepository(sql)
	viewHistoryJobRepo := history_repository.NewViewJobHistoryRepository(sql)
	viewHistoryJobUseCase := usecases.NewViewJobHistoryUseCase(viewHistoryJobRepo)
	jobRepo := jobs_repository.NewJobRepository(sql)
	cvRepo := cv_repository.NewCVRepoImpl(sql)
	// Initialize wallet repository and use case first
	walletRepo := wallet_repository.NewWalletRepository(sql)
	walletUseCase := usecases.NewWalletUseCase(walletRepo, userRepo)
	userUseCase := usecases.NewUserUseCase(
		userRepo,
		authVerifyRepo,
		&gmailService,
		funcRepo,
		roleRepo,
		uploadFileUseCase,
		aiUseCase,
		matchingJobCVRepo,
		jobRepo,
		cvRepo,
		pushService,
		walletUseCase,
	)
	userHandler := &users.UserHandler{
		UserUseCase:       userUseCase,
		CompanyUseCase:    companyUseCase,
		PermissionUseCase: permissionUseCase,
	}
	permissionHandler := &permission_handler.PermissionHandler{
		PermissionUseCase: permissionUseCase,
	}

	//Initialize Company repositories and use cases
	companyHandler := &companies.CompanyHandler{
		CompanyUseCase:    companyUseCase,
		PermissionUseCase: permissionUseCase,
		UserUseCase:       userUseCase,
		UploadFileUsecase: uploadFileUseCase,
	}

	notificationRepo := notificationrepository.NewNotificationRepository(sql)
	emailHistoryRepo := notificationrepository.NewEmailHistoryRepository(sql)
	templateRepo := notificationrepository.NewTemplateNotificationRepository(sql)
	notificationUseCase := usecases.NewNotificationUseCase(notificationRepo, emailHistoryRepo, templateRepo, &gmailService, pushService, cronTaskUseCase, userUseCase)

	userUseCase.SetNotificationUseCase(notificationUseCase)
	cvUseCase := usecases.NewCVUsecase(
		cvRepo,
		uploadFileUseCase,
		userUseCase,
		matchingJobCVRepo,
		notificationUseCase,
	)

	cvHandler := &cvs.CVHandler{
		CVUsecase:         cvUseCase,
		PermissionUseCase: permissionUseCase,
		UploadFileUseCase: uploadFileUseCase,
	}
	jobUseCase := usecases.NewJobUseCase(
		jobRepo,
		aiUseCase,
		permissionUseCase,
		uploadFileUseCase,
		thumbnailUseCase,
		userUseCase,
		matchingJobCVRepo,
		notificationUseCase,
	)

	jobHandler := &jobs.JobHandler{
		JobUseCase:            jobUseCase,
		PermissionUseCase:     permissionUseCase,
		ConfigurationUseCase:  configurationUseCase,
		UserUseCase:           userUseCase,
		ViewJobHistoryUseCase: viewHistoryJobUseCase,
		CVUseCase:             cvUseCase,
	}

	thumbnailHandler := &thumbnails.ThumbnailHandler{
		ThumbnailUseCase: thumbnailUseCase,
	}

	recruitmentCVRepo := recruitment_cv_repository.NewRecruitmentCVRepository(sql)

	recruitmentCVUseCase := usecases.NewRecruitmentCVUseCase(recruitmentCVRepo,
		permissionUseCase, notificationUseCase,
		cronTaskUseCase, uploadFileUseCase, jobUseCase,
		matchingJobCVRepo, userUseCase,
		configurationUseCase, companyUseCase)
	notificationHandler := &notification_handlers.NotificationHandler{
		NotificationUseCase: notificationUseCase,
		PermissionUseCase:   permissionUseCase,
	}
	recruitmentHandler := &recruitments.RecruitmentHandler{
		RecruitmentCVUseCase: recruitmentCVUseCase,
		JobUseCase:           jobUseCase,
		PermissionUseCase:    permissionUseCase,
		NotificationUseCase:  notificationUseCase,
		UserUseCase:          userUseCase,
		CvUseCase:            cvUseCase,
		CronTaskUseCase:      cronTaskUseCase,
		UploadUseCase:        uploadFileUseCase,
		CompanyUseCase:       companyUseCase,
	}

	// Initialize test repository and usecase
	testRepo := test_repository.NewTestRepository(sql)
	testUseCase := usecases.NewTestUseCase(testRepo)

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
