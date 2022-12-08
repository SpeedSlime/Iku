package hotupdate

type JVersion struct {
	Resource		string `json:"resVersion"`
	Client			string `json:"clientVersion"`
}

type XVersion struct {
	Platform		string  `xorm:"pk varchar(6) not null 'Platform'"`
	Resource		string	`xorm:"varchar(12) not null 'Resource'"`
	Client			string	`xorm:"varchar(12) not null 'Client'"`
}
