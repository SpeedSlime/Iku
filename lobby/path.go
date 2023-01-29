package lobby

import (
	asahi "github.com/SpeedSlime/Asahi"
	db "github.com/SpeedSlime/Asahi/database"
	"github.com/SpeedSlime/Asahi/router"
)

func Routes() []router.Route {
	asahi.Handle(db.Create(sqlLobby{}, sqlPlayer{}), "Routes")
	return []router.Route{
		router.NewGetRoute("/join", true, false, LobbyJoinGetRoute),
		router.NewPostRoute("/leave", true, false, LobbyLeavePostRoute),
		router.NewPostRoute("/ready", true, false, LobbyReadyPostRoute),
	}
}
