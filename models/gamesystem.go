package models

import "gorm.io/gorm"

type GameSystem struct {
	gorm.Model
	Name   string
	GameID uint
}
