package internal

import "time"

type Game struct {
	ID           int        `json:"id"`
	Date         *time.Time `json:"date"`
	Home         Team       `json:"home_team"`
	HomeScore    int        `json:"home_team_score"`
	Visitor      Team       `json:"visitor_team"`
	VisitorScore int        `json:"visitor_team_score"`
	Season       Season     `json:"season"`
	IsPostseason bool       `json:"postseason"`
}

// calculate winner

// calculate margin for +/-
