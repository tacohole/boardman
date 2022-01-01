package internal

import "time"

type PageData struct {
	TotalPages    int `json:"total_pages"`
	PageIndex     int `json:"current_page"`
	NextPageIndex int `json:"next_page"`
	PageSize      int `json:"per_page"`
	ItemCount     int `json:"total_count"`
}

type Page struct {
	Data     []Data   `json:"data"`
	PageData PageData `json:"meta"`
}

type Data struct {
	ID           int        `json:"id"`
	Name         string     `json:"full_name,omitempty"`
	Abbrev       string     `json:"abbreviation,omitempty"`
	Conference   string     `json:"conference,omitempty"`
	Division     string     `json:"division,omitempty"`
	FirstName    string     `json:"first_name,omitempty"`
	LastName     string     `json:"last_name,omitempty"`
	CurrentTeam  Team       `json:"team,omitempty"`
	Date         *time.Time `json:"date,omitempty"`
	Home         Team       `json:"home_team,omitempty"`
	HomeScore    int        `json:"home_team_score,omitempty"`
	Visitor      Team       `json:"visitor_team,omitempty"`
	VisitorScore int        `json:"visitor_team_score,omitempty"`
	Season       Season     `json:"season,omitempty"`
	IsPostseason bool       `json:"postseason,omitempty"`
}
