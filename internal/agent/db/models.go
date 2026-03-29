package db

import (
	"time"

	"gorm.io/gorm"
)

type Containers struct {
	gorm.Model
	ID        string
	Image     string
	Name      string
	Port      int
	Status    string
	DockerID  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
