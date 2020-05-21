/*
Package database implements the database functionality used for the chatterbot executable.
*/
package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

// Setup sets up the package global database connection used for all other functions in the package.
// Setup also ensures that the provided database has all the tables necessary for the functionality of chatterbot.
// Setup should always be called before any other function in this package.
func Setup(filepath string) error {
	var err error
	db, err = gorm.Open("sqlite3", filepath)

	if err != nil {
		return err
	}

	return nil
}

func Close() error {
	return db.Close()
}
