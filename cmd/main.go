package main

import (
	"context"
	"log"

	_ "objexrg/cmd/docs" // Import generated docs

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"objexrg/internal/config"

	auth_handler "objexrg/internal/handlers/http/auth"
	bucket_handler "objexrg/internal/handlers/http/bucket"
	storage_handler "objexrg/internal/handlers/http/storage"
	user_handler "objexrg/internal/handlers/http/user"
	"objexrg/internal/infrastructure/blob"
	"objexrg/internal/infrastructure/database"
	authrepository "objexrg/internal/infrastructure/repositories/auth"
	storage_repo "objexrg/internal/infrastructure/repositories/storage"
	"objexrg/internal/infrastructure/storage"
	"objexrg/internal/interfaces"
	"objexrg/internal/middleware"
	auth_models "objexrg/internal/models/auth"
	storage_models "objexrg/internal/models/storage"

	authsvc "objexrg/internal/services/auth"
	bucketsvc "objexrg/internal/services/bucket"
	services "objexrg/internal/services/storage"
	usersvc "objexrg/internal/services/user"
	authval "objexrg/internal/validators/auth"
)

// @title ObjexRG Backend API
// @version 1.0
// @description A simple S3-like storage API
// @host localhost:7080
// @BasePath /

// 🔐 Swagger Bearer config
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.LoadConfig()
	basePath := cfg.BasePath

	db := database.Connect()

	db.AutoMigrate(
		&auth_models.UserData{}, &auth_models.UserCredentials{}, &auth_models.RefreshToken{},
		&storage_models.Bucket{}, &storage_models.Objects{}, &storage_models.ObjectShares{},
	)

	storageRepo := storage_repo.NewStorageRepository(db)

	blobStore, err := blob.NewBlobStore(context.Background(), cfg)
	if err != nil {
		panic("Failed to initialize storage: " + err.Error())
	}

	var bucketBackend interfaces.BucketInterface
	if cfg.StorageDriver == blob.DriverS3 {
		bucketBackend = storage.NewNoopBucketStorage()
	} else {
		bucketBackend = storage.NewLocalStorage(basePath)
	}

	bucketService := bucketsvc.NewBucketService(bucketBackend, storageRepo)
	uploadService := services.NewUploadService(blobStore, storageRepo, cfg.MaxUploadBytes, cfg.PresignTTL)
	shareService := services.NewShareService(storageRepo)
	objectService := services.NewObjectService(blobStore, storageRepo, cfg.PresignTTL)

	delivery := storage_handler.NewDelivery(blobStore, cfg.PresignTTL)

	authRepo := authrepository.NewAuthRepository(db)
	tokenService, err := authsvc.NewTokenService()
	if err != nil {
		panic("Failed to initialize token service: " + err.Error())
	}
	userService := usersvc.NewUserService(authRepo, db)
	authService := authsvc.NewAuthService(authRepo, db, tokenService, bucketService)
	registerValidator := authval.NewRegisterValidator(authRepo)
	loginValidator := authval.NewLoginValidator(authRepo)
	authHandler := auth_handler.NewAuthHandler(authService, registerValidator, loginValidator)
	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authGroup := router.Group("/auth")
	auth_handler.RegisterAuthRoutes(authGroup, authHandler)

	// Protected routes
	protected := router.Group("/")
	protected.Use(authMiddleware.JWTAuth())
	{
		userGroup := protected.Group("/user")
		user_handler.RegisterUserRoutes(userGroup, user_handler.NewUserHandler(userService))

		bucketGroup := protected.Group("/bucket")
		bucket_handler.RegisterBucketRoutes(bucketGroup, bucketService)

		uploadGroup := protected.Group("/upload")
		storage_handler.RegisterUploadRoutes(uploadGroup, storage_handler.NewUploadHandler(uploadService, cfg.MaxUploadBytes))

		downloadGroup := protected.Group("/download")
		storage_handler.RegisterDownloadRoutes(downloadGroup, storage_handler.NewDownloadHandler(objectService, delivery))

		storage_handler.RegisterObjectRoutes(protected, storage_handler.NewObjectHandler(objectService, delivery, cfg.PublicBaseURL))

		storage_handler.RegisterShareRoutes(protected, router.Group("/"), storage_handler.NewShareHandler(shareService, delivery, cfg.PublicBaseURL))
		storage_handler.RegisterPublicObjectRoutes(router.Group("/"), storage_handler.NewObjectHandler(objectService, delivery, cfg.PublicBaseURL))
	}

	if err := router.Run(":7080"); err != nil {
		panic(err)
	}
}
