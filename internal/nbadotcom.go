package internal

// nba.com endpoint consts
const (
	NbaDataUrl = "http://data.nba.net/prod/v1/"
	Coaches    = "/coaches.json"
	Teams      = "/teams.json"
)

// structs
type NbaPage struct {
	Internal []struct{} `json:"internal"` // we don't care about this
	League   NbaLeague  `json:"league"`
}

type NbaLeague struct {
	Standard []NbaData `json:"standard"`
}

type NbaData struct {
	Name        string `json:"fullName"`
	Abbrev      string `json:"tricode"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	IsAssistant bool   `json:"isAssistant"`
	PersonID    string `json:"personId"`
	TeamID      string `json:"teamId"`
}
