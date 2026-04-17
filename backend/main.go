package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"nav/backend/config"
	"nav/backend/database"
	"nav/backend/models"
	"nav/backend/routes"
	"nav/backend/utils"

	"gorm.io/gorm"
)

func init() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}
	time.Local = loc
}

func main() {
	config.LoadConfig()

	database.InitDB()
	db := database.GetDB()

	models.Migrate(db)
	ensureDefaultAdmin(db)

	router := routes.SetupRouter()

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func ensureDefaultAdmin(db *gorm.DB) {
	var user models.User
	if err := db.Where("username = ?", "admin").First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Fatalf("Failed to query default admin user: %v", err)
		}

		hashedPassword, err := utils.HashPassword("admin")
		if err != nil {
			log.Fatalf("Failed to hash default admin password: %v", err)
		}

		now := time.Now()
		user = models.User{
			Username: "admin",
			Password: hashedPassword,
			LoginAt:  models.NewNullableTime(now),
			Status:   "enabled",
			IsAdmin:  true,
		}

		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("Failed to create default admin user: %v", err)
		}
		log.Println("Initialized default admin user with username: admin and password: admin")
		return
	}

	if user.IsAdmin {
		log.Println("Admin user already exists and is an administrator")
		return
	}

	user.IsAdmin = true
	if err := db.Save(&user).Error; err != nil {
		log.Fatalf("Failed to update admin user to administrator: %v", err)
	}
	log.Println("Updated existing admin user to have administrator privileges")
}
