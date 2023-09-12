package models

import (
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	UserID      uint // Use the user's ID as the foreign key
	Status      string
	Description string
}
