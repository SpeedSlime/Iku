package leaderboard

// RecordID is a snowflake of the lobbyid + unix timestamp + has of the record there is a very high 

type sqlLeaderboard struct {
	RecordID		string	`xorm:"pk varchar(12) not null 'RecordID'"`
	Name			string	`xorm:"varchar(12) not null 'Name'"`
	Seed			string	`xorm:"varchar(12) 'Seed'"`
	Time			int64	`xorm:"bigint notnull 'Time'"`
}