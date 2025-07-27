package db

import (
	"fmt"
	"log"
	"mini-crm/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	dsn := config.GetDatabaseUrl()
	fmt.Println("DSN:", dsn)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
}
