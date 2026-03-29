package main

import (
	"log"

	_ "objexrg-backend/cmd/docs" // Import generated docs

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"objexrg-backend/internal/config"

	auth_handler "objexrg-backend/internal/handlers/http/auth"
	bucket_handler "objexrg-backend/internal/handlers/http/bucket"
	storage_handler "objexrg-backend/internal/handlers/http/storage"
	user_handler "objexrg-backend/internal/handlers/http/user"
	"objexrg-backend/internal/infrastructure/database"
	authrepository "objexrg-backend/internal/infrastructure/repositories/auth"
	"objexrg-backend/internal/infrastructure/storage"
	"objexrg-backend/internal/middleware"
	models "objexrg-backend/internal/models/auth"

	authsvc "objexrg-backend/internal/services/auth"
	bucketsvc "objexrg-backend/internal/services/bucket"
	services "objexrg-backend/internal/services/storage"
	usersvc "objexrg-backend/internal/services/user"
	authval "objexrg-backend/internal/validators/auth"
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

	db.AutoMigrate(&models.UserData{}, &models.UserCredentials{})

	localStore := storage.NewLocalStorage(basePath)
	bucketService := bucketsvc.NewBucketService(localStore)
	uploadService := services.NewUploadService(localStore)

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
		bucket_handler.RegisterBucketRoutes(bucketGroup, bucketService, localStore)

		uploadGroup := protected.Group("/upload")
		storage_handler.RegisterUploadRoutes(uploadGroup, storage_handler.NewUploadHandler(uploadService))

		downloadGroup := protected.Group("/download")
		storage_handler.RegisterDownloadRoutes(downloadGroup, uploadService)
	}

	if err := router.Run(":7080"); err != nil {
		panic(err)
	}
}
