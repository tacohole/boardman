package schema

type PageData struct {
	TotalPages    int `json:"total_pages"`
	PageIndex     int `json:"current_page"`
	NextPageIndex int `json:"next_page"`
	PageSize      int `json:"per_page"`
	ItemCount     int `json:"total_count"`
}

type Player struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	CurrentTeam Team   `json:"team"`
}

type Team struct {
	ID         int    `json:"id"`
	Name       string `json:"full_name"`
	Abbrev     string `json:"abbreviation"`
	Conference string `json:"conference"`
	Division   string `json:"division"`
}

type Season struct {
	LeagueYear string
	Champion   Team
	WConfChamp Team
	EConfChamp Team
	MVP        Player
}

type TeamYear struct {
	TeamCache    Team
	Season       Season
	Wins         int
	Losses       int
	WinPct       int
	ConfRank     int
	OvrRank      int
	MadePlayoffs bool
	Roster       []Player
	Coach        string
}

type PlayerYear struct {
	Player Player
	Season Season
	Stats  []int
}
