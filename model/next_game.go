package model

import "time"

type NextGame struct {
	MatchID    int
	HomeTeam   string
	HomeTeamID int
	AwayTeam   string
	AwayTeamID int
	StartTime  time.Time
}
