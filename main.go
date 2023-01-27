package main

import (
	"github.com/SpeedSlime/Iku/lobby"

	"github.com/SpeedSlime/Asahi"
	"github.com/SpeedSlime/Asahi/middleware"
	"github.com/SpeedSlime/Asahi/router"

	alternate "github.com/go-chi/chi/middleware"
)

func main() {
	s := asahi.New()
	s.LoadMiddleware(Middlewares())
	s.LoadRouter(Routers())
	s.Start()
}

func Middlewares() []middleware.Middleware {
	return []middleware.Middleware{
		middleware.NewMiddleware(alternate.Logger, true, false),
	}
}

func Routers() []router.Router {
	return []router.Router{
		router.NewRouter(lobby.Routes(), true),
	}
}