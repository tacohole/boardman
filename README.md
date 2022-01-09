# boardman gets paid

## building a dataset of NBA team/player stats from [balldontlie.io](https://www.balldontlie.io/)

## Pre-Reqs
- go v1.17
- PostgreSQL

## Install and Run
### local installation
`mkdir boardman`
`cd boardman`
`gh repo clone tacohole/boardman`

### configuration
`cat boardman-config.env`
- add your database URL as "DATABASE_URL=postgresql://user:secret@host:port/database_name"
- set your database timeout as "DB_TIMEOUT=30s" (accepts ms/s/m/h)

### compile and run
`go build -o boardman`

#### get teams
`./boardman get teams`

#### get players
`./boardman get players`

#### get games
`./boardman get games`

#### get player season averages
`./boardman get player-stats`

#### get player single game stats
`./boardman get game-stats`

#### get coaches
`./boardman get coaches`

##### get champions and award winners
`./boardman get champs`

### get all
`./boardman get paid`