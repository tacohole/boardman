# boardman gets paid

## building a dataset of NBA team/player stats from [balldontlie.io](https://www.balldontlie.io/) and [data.nba.com](https://data.nba.com/)

## Pre-Reqs
- go >=v1.14
- PostgreSQL >=v13

## Install and Run

### local installation
```
mkdir boardman
cd boardman
gh repo clone tacohole/boardman
```

### configuration
`cat boardman-config.env`
- add your database URL as `DATABASE_URL=postgresql://user:secret@host:port/database_name`
- set your database timeout as `DB_TIMEOUT=30s` (accepts ms/s/m/h)

### compile
`go build -o boardman`

### what it does:
each command creates the relevant relations in Posgres
and populates tables with data

#### get teams
`./boardman get teams`

#### get players
`./boardman get players`

#### get coaches
`./boardman get coaches`

#### get games
`./boardman get games`

#### get player season averages
`./boardman get player-stats`

#### get player single game stats
`./boardman get game-stats`

#### get all
`./boardman get paid`

shout out to [nbasense.com](http://nbasense.com/) for NBA API documentation