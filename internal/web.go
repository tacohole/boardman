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
}
