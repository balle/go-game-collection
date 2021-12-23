package models

import (
	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	Name         string
	Played       bool
	Finished     bool
	GameSystemID int
	GameSystems  []GameSystem
}
