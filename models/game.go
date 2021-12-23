package models

import (
	"time"

	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	Name         string
	Played       bool
	Finished     bool
	GameSystemID int
	GameSystems  []GameSystem
	StartedAt    time.Time
	FinishedAt   time.Time
}
