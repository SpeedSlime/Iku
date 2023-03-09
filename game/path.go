package game

import (
	asahi "github.com/SpeedSlime/Asahi"
	db "github.com/SpeedSlime/Asahi/database"
	"github.com/SpeedSlime/Asahi/router"
)

func Routes() []router.Route {
	asahi.Handle(db.Create(Maze{}), "Routes")
	return []router.Route{
		router.NewPostRoute("/getMaze", true, false, GameSeedPostRoute),
	}
}