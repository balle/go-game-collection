package store

import (
	"github.com/balle/go-game-collection/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

/*
 * VARIABLES
 */

var Db *gorm.DB

func init() {
	// Connect to the database
	var err error
	Db, err = gorm.Open(sqlite.Open("games.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	Db.AutoMigrate(&models.GameSystem{}, &models.Game{})

}
