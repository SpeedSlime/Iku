package game

type Maze struct {
	Width  int       `json:"width"`
	Height int       `json:"height"`
	Grid   [][]*Cell `json:"grid"`
}

type Cell struct {
	Row     int     `json:"row"`
	Column  int     `json:"column"`
	Walls   [4]bool `json:"walls"`
	Visited bool    `json:"visited"`
}