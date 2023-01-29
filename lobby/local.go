package lobby

import (
	"net/http"

	"github.com/SpeedSlime/Asahi"
	db "github.com/SpeedSlime/Asahi/database"
	"github.com/SpeedSlime/Asahi/reply"
)

func LobbyJoinGetRoute(w http.ResponseWriter, r *http.Request) {
	var players []sqlPlayer
	var lobby jsonLobby

	name := r.URL.Query().Get("name")
	texture := r.URL.Query().Get("texture")
	code := r.URL.Query().Get("code")

	p := sqlPlayer{Name: name, LobbyID: code}
	// We are searching for duplicate names, so we cannot include texture in the above.
	if p.Name == "" || texture == "" || code == "" || len(code) > 6 || len(p.Name) > 12 || len(texture) > 32 {
		// In this case, it is the clients fault for sending a bad request.
		asahi.Handle(reply.RespondWithResult(w, http.StatusForbidden, "You are not allowed to do that."), "LobbyJoinGetRoute"); return
	}

	if db.Exists(&p) {
		// If the user exists, we need to throw an error.
		asahi.Handle(reply.RespondWithResult(w, http.StatusConflict, "A user in the specified lobby already has that name!"), "LobbyJoinGetRoute"); return
	}

	// No user in the specified lobby exists, thus we need to check that the lobby actually exists.
	l := sqlLobby{LobbyID: code}
	if has, _ := db.Select(&l); !has {
		// Lobby does not exist
		asahi.Handle(reply.RespondWithResult(w, http.StatusBadRequest, "Lobby not found."), "LobbyJoinGetRoute"); return
	}

	// If the user has reached this, it means that a lobby exists and there is a free username.
	// We need to also check if a game is in progress, defined below, we check if InProgress is true or if server is at max player-count.
	if l.InProgress {
		// User cannot join session
		asahi.Handle(reply.RespondWithResult(w, http.StatusConflict, "The lobby is already in progress."), "LobbyJoinGetRoute"); return
	}

	// We now need to fetch any users in the lobby. We don't care if there are no players in the lobby, as the server periodically checks.
	_, err := db.Find(&players, &sqlPlayer{LobbyID: code})
	if int(l.MaxPlayers) <= len(players) {
		asahi.Handle(reply.RespondWithResult(w, http.StatusConflict, "The server is full."), "LobbyJoinGetRoute"); return
	}

	if err != nil {
		// This will be an internal server error. User fault isn't possible here.
		asahi.Handle(reply.RespondWithResult(w, http.StatusInternalServerError, "An unexpected error has occured when trying to find the server, please try again later."), "LobbyJoinGetRoute"); return
	}

	// We then insert the user.
	p.Points = 0
	p.PreviousTime = 0
	p.Rank = -1
	p.Ready = false
	p.Texture = texture
	err = db.Insert(&p)
	if err != nil {
		// This would mean the server is at fault, and is an unexpected error.
		asahi.Handle(reply.RespondWithResult(w, http.StatusInternalServerError, "An unexpected error has occured when trying to join the server, please try again later."), "LobbyJoinGetRoute"); return
	}

	// Construct the json.
	for _, player := range players {
		buff :=	jsonPlayer{Name: player.Name, Texture: player.Texture, Rank: player.Rank, Points: player.Points, PreviousTime: player.PreviousTime, Ready: player.Ready}
		lobby.Players = append(lobby.Players, buff)
	}
	lobby.EndRound = l.EndRound
	lobby.LobbyCode = l.LobbyID
	lobby.MaxPlayers = l.MaxPlayers
	lobby.Round = l.Round
	lobby.Host = l.Host
	lobby.Seed = l.Seed

	// User can join session, return details.
	asahi.Handle(reply.RespondWithJSON(w, http.StatusOK, lobby), "LobbyJoinGetRoute")
}

func LobbyReadyPostRoute(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.FormValue("code")
	name := r.FormValue("name")

	p := sqlPlayer{Name: name, LobbyID: code}

	if p.LobbyID == "" || p.Name == "" {
		// Check for cases that have malformed requests.
		asahi.Handle(reply.RespondWithResult(w, http.StatusForbidden, "You are not allowed to do that."), "LobbyReadyPostRoute"); return
	}

	// Make a copy of the params and set ready to true
	n := p
	n.Ready = true
	p.Ready = false

	// We only want to set the "Ready" collumn to true, so we pass the "p" as the cond
	has, err := db.Update(&n, &p)
	
	if err != nil {
		// If there is an error here, it will be due to a server error.
		asahi.Handle(reply.RespondWithResult(w, http.StatusInternalServerError, "An unexpected error has occured when readying, you have kicked from the session."), "LobbyReadyPostRoute"); return
	}

	if !has {
		// If the user doesn't exist, we need to throw an error.
		asahi.Handle(reply.RespondWithResult(w, http.StatusConflict, "There has been a desynchronisation between the client and server."), "LobbyReadyPostRoute"); return
	}

	// We do not need to handle any other logic here in terms of starting the match, as that is handled by the game host.
	asahi.Handle(reply.RespondWithResult(w, http.StatusOK, "Client is Ready."), "LobbyReadyPostRoute")
}


func LobbyLeavePostRoute(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	name := r.FormValue("name")

	p := sqlPlayer{Name: name, LobbyID: code}

	if !db.Exists(&p) {
		// It's possible that multiple requests could be sent in the case of network disconnects, so, we need to notify the client has already left just in case
		asahi.Handle(reply.RespondWithResult(w, http.StatusOK, "Client has successfully left."), "LobbyLeavePostRoute"); return
	}

	if err := db.Delete(&p); err != nil {
		// This is a fatal server error. It *must* be handled.
		asahi.Handle(reply.RespondWithResult(w, http.StatusInternalServerError, "An unexpected error has occured."), "LobbyLeavePostRoute"); return
	}

	asahi.Handle(reply.RespondWithResult(w, http.StatusOK, "Client has successfully left."), "LobbyLeavePostRoute")
}