package lobby

// In this case, all data here is temporary, so no need to store in a seperate table, as this is also fetched on server join.
// Only data that needs to be stored seperately is round information.
// Also, We must check no duplicate names!

// type sqlLeaderboard struct {
// 	RecordID		string	`xorm:"pk varchar(12) not null 'RecordID'"`
// 	Name			string	`xorm:"varchar(12) not null 'Name'"`
// 	Seed			string	`xorm:"varchar(12) 'Seed'"`
// 	Time			int64	`xorm:"bigint notnull 'Time'"`
// }

type sqlPlayer struct {
	LobbyID			string	`xorm:"pk varchar(6) not null 'LobbyID'"`
	Name			string	`xorm:"pk varchar(12) not null 'Name'"`
	Texture			string	`xorm:"varchar(32) not null 'TexturePath'"`
	Rank			int64	`xorm:"bigint not null 'Rank'"`
	Points			int64	`xorm:"bigint not null 'Points'"`
	PreviousTime 	int64	`xorm:"bigint not null 'PreviousTime'"`
}

type sqlLobby struct {
	LobbyID			string 	`xorm:"pk varchar(6) not null 'LobbyID'"`
	Host			string	`xorm:"varchar(12) not null 'Host'"`		
	EndRound		int64	`xorm:"bigint not null 'EndRound'"`
	MaxPlayers		int64	`xorm:"bigint not null 'MaxPlayer'"`
	Round			int64	`xorm:"bigint not null 'Round'"`
	Seed			string	`xorm:"varchar(12) 'Seed'"`
	InProgress		bool	`xorm:"bool 'InProgress'"`
}

type jsonPlayer struct {
	Name			string	`json:"nickname"`
	Texture			string	`json:"skin"`
	Rank			int64	`json:"rank"`
	Points			int64	`json:"points"`
	PreviousTime	int64	`json:"ts"`
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