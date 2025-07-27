package main

import (
	"mini-crm/internal/config"
	"mini-crm/internal/db"
	"mini-crm/internal/handler"
	"mini-crm/internal/middleware"
	"mini-crm/internal/repositoryimpl"
	"mini-crm/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	db.Init()

	r := gin.Default()

	// Миграция
	db.DB.AutoMigrate()

	// DI: Сборка зависимостей
	userRepo := repositoryimpl.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Роуты
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.GET("/ping", middleware.AuthMiddleware(userRepo), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.Run()
}
