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
	code := asahi.Parameter(r, "code")

	p := sqlPlayer{Name: name, LobbyID: code}
	// We are searching for duplicate names, so we cannot include texture in the above.
	if p.Name == "" || texture == "" || len(p.Name) > 12 || len(texture) > 32 {
		// In this case, it is the clients fault for sending a bad request.
		asahi.Handle(reply.RespondWithResult(w, http.StatusBadRequest, "Please provide a valid nickname and texture."), "LobbyJoinGetRoute"); return
	}

	// We don't want to pass by reference here, as it could make an unnessesary overwrite.
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
	p.Texture = texture
	err = db.Insert(&p)
	if err != nil {
		// This would mean the server is at fault, and is an unexpected error.
		asahi.Handle(reply.RespondWithResult(w, http.StatusInternalServerError, "An unexpected error has occured when trying to join the server, please try again later."), "LobbyJoinGetRoute"); return
	}

	// Construct the json.
	for _, player := range players {
		buff :=	jsonPlayer{Name: player.Name, Texture: player.Texture, Rank: player.Rank, Points: player.Points, PreviousTime: player.PreviousTime}
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

