package main

import (
	"fmt"
	"log"
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

	if err := db.CheckConnection(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	fmt.Println("Database connection verified")

	r := gin.Default()

	userRepo := repositoryimpl.NewUserRepository(db.DB)
	tokenRepo := repositoryimpl.NewTokenRepository(db.DB)
	userService := service.NewUserService(userRepo, tokenRepo)
	userHandler := handler.NewUserHandler(userService)

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.POST("/refresh", userHandler.Refresh)
	r.POST("/logout", userHandler.Logout)
	r.GET("/me", middleware.AuthMiddleware(userRepo), userHandler.Me)

	r.Run()
}
