package main

import (
	"golang-rest-api-jwt/config"
	v1 "golang-rest-api-jwt/handler/v1"
	"golang-rest-api-jwt/repository"
	"golang-rest-api-jwt/service"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.SetupDatabaseConnection()
	defer config.CloseDatabaseConnection(db)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	authService := service.NewAuthService(userRepository)
	jwtService := service.NewJWTService()
	authHandler := v1.NewAuthHandler(userService, authService, jwtService)

	server := gin.Default()
	authRoutes := server.Group("api/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	server.Run()
}
