package auth

import (
	asahi "github.com/SpeedSlime/Asahi"
	db "github.com/SpeedSlime/Asahi/database"
	"github.com/SpeedSlime/Asahi/router"
)

func Routes() []router.Route {
	asahi.Handle(db.Create(User{}), "Routes")
	return []router.Route{
		router.NewPostRoute("/login", true, false, AuthLoginPostRoute),
		router.NewPostRoute("/create", true, false, AuthCreatePostRoute),
		router.NewPostRoute("/logout", true, false, AuthLogoutPostRoute),
	}
}