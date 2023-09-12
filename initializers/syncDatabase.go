package initializers

import "examples.com/jwt-auth/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Todo{})
	DB.AutoMigrate(&models.BlackListedToken{})
}
