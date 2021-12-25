package internal

type Season struct {
	LeagueYear string
	Champion   Team
	WConfChamp Team
	EConfChamp Team
	MVP        Player
}

// get champ
// get conf champ
// get mvp

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

// sum wins
// sum losses
// sum wpct
// get conf rank
// build roster

type PlayerYear struct {
	Player Player
	Season Season
	Stats  []int
}
