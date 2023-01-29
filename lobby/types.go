package lobby

type sqlPlayer struct {
	LobbyID			string	`xorm:"pk varchar(6) not null 'LobbyID'"`
	Name			string	`xorm:"pk varchar(12) not null 'Name'"`
	Texture			string	`xorm:"varchar(32) not null 'TexturePath'"`
	Rank			int64	`xorm:"bigint not null 'Rank'"`
	Points			int64	`xorm:"bigint not null 'Points'"`
	Ready			bool	`xorm:"bool not null 'Ready'"`
	PreviousTime 	int64	`xorm:"bigint not null 'PreviousTime'"`
}

type sqlLobby struct {
	LobbyID			string 	`xorm:"pk varchar(6) not null 'LobbyID'"`
	Host			string	`xorm:"varchar(12) not null 'Host'"`		
	EndRound		int64	`xorm:"bigint not null 'EndRound'"`
	MaxPlayers		int64	`xorm:"bigint not null 'MaxPlayer'"`
	Round			int64	`xorm:"bigint not null 'Round'"`
	Seed			string	`xorm:"varchar(12) 'Seed'"`
	InProgress		bool	`xorm:"bool not null 'InProgress'"`
}

type jsonPlayer struct {
	Name			string	`json:"nickname"`
	Texture			string	`json:"skin"`
	Rank			int64	`json:"rank"`
	Points			int64	`json:"points"`
	PreviousTime	int64	`json:"ts"`
	Ready			bool	`json:"isReady"`
}

type jsonLobby struct {
	Players		[]jsonPlayer	`json:"players"`
	Round		int64			`json:"round"`
	EndRound	int64			`json:"endRound"`
	MaxPlayers	int64			`json:"maxPlayers"`
	LobbyCode	string			`json:"code"`
	Host		string			`json:"host"`
	Seed		string			`json:"seed"`
}