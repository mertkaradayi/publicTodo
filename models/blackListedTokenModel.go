package models

import (
	"time"

	"gorm.io/gorm"
)

type BlackListedToken struct {
	gorm.Model
	Token     string
	ExpiresAt time.Time
}
