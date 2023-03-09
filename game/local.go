package game

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/SpeedSlime/Iku/auth"

	asahi "github.com/SpeedSlime/Asahi"
	db "github.com/SpeedSlime/Asahi/database"
	"github.com/SpeedSlime/Asahi/reply"
)

const (
	width = 50
	height = 50
)

func GenerateMaze(width, height int, seed int64) *Maze {
	rand.Seed(seed)
	maze := &Maze{Width: width, Height: height}
	grid := make([][]*Cell, height)
	for i := range grid {
		grid[i] = make([]*Cell, width)
		for j := range grid[i] {
			grid[i][j] = &Cell{Row: i, Column: j, Walls: [4]bool{true, true, true, true}, Visited: false}
		}
	}

	current := grid[0][0]
	current.Visited = true

	stack := []*Cell{current}

	for len(stack) > 0 {
		index := rand.Intn(len(stack))
		current := stack[index]

		neighbors := []*Cell{}
		if current.Row > 0 {
			neighbors = append(neighbors, grid[current.Row-1][current.Column])
		}
		if current.Row < height-1 {
			neighbors = append(neighbors, grid[current.Row+1][current.Column])
		}
		if current.Column > 0 {
			neighbors = append(neighbors, grid[current.Row][current.Column-1])
		}
		if current.Column < width-1 {
			neighbors = append(neighbors, grid[current.Row][current.Column+1])
		}

		unvisitedNeighbors := []*Cell{}
		for _, neighbor := range neighbors {
			if !neighbor.Visited {
				unvisitedNeighbors = append(unvisitedNeighbors, neighbor)
			}
		}

		if len(unvisitedNeighbors) > 0 {
			neighborIndex := rand.Intn(len(unvisitedNeighbors))
			neighbor := unvisitedNeighbors[neighborIndex]

			if neighbor.Row < current.Row {
				current.Walls[0] = false
				neighbor.Walls[1] = false
			} else if neighbor.Row > current.Row {
				current.Walls[1] = false
				neighbor.Walls[0] = false
			} else if neighbor.Column < current.Column {
				current.Walls[2] = false
				neighbor.Walls[3] = false
			} else if neighbor.Column > current.Column {
				current.Walls[3] = false
				neighbor.Walls[2] = false
			}

			neighbor.Visited = true
			stack = append(stack, neighbor)
		} else {
			stack = append(stack[:index], stack[index+1:]...)
		}
	}

	maze.Grid = grid
	return maze
}

func GameSeedPostRoute(w http.ResponseWriter, r *http.Request) {
	var user auth.User

	r.ParseForm()
	seed := r.FormValue("seed")
	token := r.FormValue("token")
	user.Username = r.FormValue("username")

	if token == "" || user.Username == "" || seed == "" {
		asahi.Handle(reply.RespondWithResult(w, http.StatusForbidden, "A field is misssing."), "GameSeedPostRoute"); return
	}

    found, err := db.Select(&user)
    if err != nil {
        asahi.Handle(reply.RespondWithResult(w, http.StatusInternalServerError, "An unexpected error has occcured."), "GameSeedPostRoute"); return 
    }

    if !found {
        asahi.Handle(reply.RespondWithResult(w, http.StatusUnauthorized, "That account does not exist."), "GameSeedPostRoute"); return  
    }

	if !(user.Token == token) {
		asahi.Handle(reply.RespondWithResult(w, http.StatusUnauthorized, "Token is invalid."), "GameSeedPostRoute"); return  
	}

	if auth.IsTokenExpired(user.TokenTime) {
		asahi.Handle(reply.RespondWithResult(w, http.StatusUnauthorized, "Token has expired."), "GameSeedPostRoute"); return  
	}

	i, err := strconv.ParseInt(seed, 10, 64)
	if err != nil {
		asahi.Handle(reply.RespondWithResult(w, http.StatusInternalServerError, "An unexpected error has occcured."), "GameSeedPostRoute"); return 
	}

	asahi.Handle(reply.RespondWithJSON(w, http.StatusOK, GenerateMaze(width, height, i)), "GameSeedPostRoute")
}