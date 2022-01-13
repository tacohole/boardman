package internal

// nba.com endpoint consts
const (
	NbaDataUrl = "http://data.nba.net/prod/v1/"
	Coaches    = "/coaches.json"
	Champs     = "/playoffsBracket.json"
	Teams      = "/teams.json"
)

// structs
type NbaPage struct {
	Internal []struct{} `json:"internal"`
	League   NbaLeague  `json:"league"`
}

type NbaLeague struct {
	Standard []TeamResponse `json:"standard"`
}

// we really just want the teamID and name here
type TeamResponse struct {
	Name   string `json:"fullName"`
	ID     string `json:"teamId"`
	Abbrev string `json:"tricode"`
}

type CoachResponse struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	IsAssistant bool   `json:"isAssistant"`
	PersonID    string `json:"personId"`
	TeamID      string `json:"teamId"`
}

type ChampResponse struct {
}
