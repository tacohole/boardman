package internal

import "time"

// PageData is page metadata for balldontlie API responses
type PageData struct {
	TotalPages    int `json:"total_pages"`
	PageIndex     int `json:"current_page"`
	NextPageIndex int `json:"next_page"`
	PageSize      int `json:"per_page"`
	ItemCount     int `json:"total_count"`
}

// Page is a page of response from balldontlie API
type Page struct {
	Data     []Data   `json:"data"`
	PageData PageData `json:"meta"`
}

// Data contains all possible fields from balldontlie API endpoints
type Data struct {
	ID           int        `json:"id"`
	Name         string     `json:"full_name"`
	Abbrev       string     `json:"abbreviation"`
	Conference   string     `json:"conference"`
	Division     string     `json:"division"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	CurrentTeam  Team       `json:"team"`
	Date         *time.Time `json:"date"`
	Home         Team       `json:"home_team"`
	HomeScore    int        `json:"home_team_score"`
	Visitor      Team       `json:"visitor_team"`
	VisitorScore int        `json:"visitor_team_score"`
	LeagueYear   int        `json:"season"`
	IsPostseason bool       `json:"postseason"`
	Minutes      string     `json:"avg_min" db:"avg_min"`
	FGM          float32    `json:"fgm" db:"fgm"`
	FGA          float32    `json:"fga" db:"fga"`
	FG3M         float32    `json:"fg3m" db:"fg3m"`
	FG3A         float32    `json:"fg3a" db:"fg3a"`
	OREB         float32    `json:"oreb" db:"oreb"`
	DREB         float32    `json:"dreb" db:"dreb"`
	REB          float32    `json:"reb" db:"reb"`
	AST          float32    `json:"ast" db:"ast"`
	STL          float32    `json:"stl" db:"stl"`
	BLK          float32    `json:"blk" db:"blk"`
	TO           float32    `json:"to" db:"to"`
	PF           float32    `json:"pf" db:"pf"`
	PTS          float32    `json:"pts" db:"pts"`
	FG_PCT       float32    `json:"fg_pct" db:"fg_pct"`
	FG3_PCT      float32    `json:"fg3_pct" db:"fg3_pct"`
	FT_PCT       float32    `json:"ft_pct" db:"ft_pct"`
}
