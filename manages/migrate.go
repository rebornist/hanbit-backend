package manages

import (
	"github.com/rebornist/hanbit/boards"
	"github.com/rebornist/hanbit/config"
	"github.com/rebornist/hanbit/gallaries"
	"github.com/rebornist/hanbit/images"
	"github.com/rebornist/hanbit/sermons"
	"github.com/rebornist/hanbit/users"
)

func Migrate() error {
	// db connect
	db := config.ConnectDb()

	board := new(boards.Board)
	gallary := new(gallaries.Gallary)
	image := new(images.Image)
	category := new(images.Category)
	sermon := new(sermons.Sermon)
	user := new(users.User)

	// Migrate the schema
	if err := db.AutoMigrate(&user); err != nil {
		return err
	}

	if err := db.AutoMigrate(&sermon); err != nil {
		return err
	}

	if err := db.AutoMigrate(&image); err != nil {
		return err
	}

	if err := db.AutoMigrate(&category); err != nil {
		return err
	}

	if err := db.AutoMigrate(&board); err != nil {
		return err
	}

	if err := db.AutoMigrate(&gallary); err != nil {
		return err
	}

	return nil
}
