package manages

import (
	"github.com/rebornist/hanbit/boards"
	"github.com/rebornist/hanbit/config"
	"github.com/rebornist/hanbit/gallaries"
	"github.com/rebornist/hanbit/images"
	"github.com/rebornist/hanbit/sermonbroadcasts"
	"github.com/rebornist/hanbit/sermons"
	"github.com/rebornist/hanbit/users"
)

func Migrate() error {
	// db connect
	db := config.ConnectDb()

	board := new(boards.Board)
	gallary := new(gallaries.Gallary)
	image := new(images.Image)
	sermon := new(sermons.Sermon)
	user := new(users.User)
	post := new(config.Post)
	logger := new(config.Logger)
	broadcast := new(sermonbroadcasts.Broadcast)

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

	if err := db.AutoMigrate(&board); err != nil {
		return err
	}

	if err := db.AutoMigrate(&gallary); err != nil {
		return err
	}

	if err := db.AutoMigrate(&post); err != nil {
		return err
	}

	if err := db.AutoMigrate(&logger); err != nil {
		return err
	}

	if err := db.AutoMigrate(&broadcast); err != nil {
		return err
	}

	return nil
}
