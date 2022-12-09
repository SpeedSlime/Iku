package hotupdate

import (
	asahi "github.com/SpeedSlime/Asahi"
	db "github.com/SpeedSlime/Asahi/database"
	"github.com/SpeedSlime/Asahi/router"
)

func Routes() []router.Route {
	asahi.Handle(db.Create(XVersion{}), "Routes")
	return []router.Route{
		router.NewGetRoute("/hotupdate/{device}/version", true, false, HotUpdateVersionGetRoute),
		router.NewGetRoute("/hotupdate/{device}/assets/{file}", true, false, HotUpdateDownloadGetRoute),
	}
}
