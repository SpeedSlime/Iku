package hotupdate

import (
	db "github.com/SpeedSlime/Asahi/database"
	"github.com/SpeedSlime/Asahi/router"
)

func Routes() []router.Route {
	db.Create(&XVersion{})
	return []router.Route{
		router.NewGetRoute("/hotupdate/{device}/version", true, false, HotUpdateVersionGetRoute),
		router.NewGetRoute("/hotupdate/{device}/assets/{file}", true, false, HotUpdateDownloadGetRoute),
	}
}
