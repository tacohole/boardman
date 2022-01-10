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
	Data []struct{} `json:"standard"`
}

// we really just want the teamID and name here
type TeamResponse struct {
	Name string `json:"fullName"`
	ID   string `json:"teamId"`
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
