package lobby

import (
	asahi "github.com/SpeedSlime/Asahi"
	db "github.com/SpeedSlime/Asahi/database"
	"github.com/SpeedSlime/Asahi/router"
)

func Routes() []router.Route {
	asahi.Handle(db.Create(sqlLobby{}, sqlPlayer{}), "Routes")
	return []router.Route{
		router.NewGetRoute("/join/{code}", true, false, LobbyJoinGetRoute),
	}
}
