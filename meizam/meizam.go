package meizam

import (
	"database/sql"
	"fmt"
	"github.com/meizam-bot/model"
)

type Meizam struct {
	connString string
	db         *sql.DB
}

func NewMeizam(connString string) *Meizam {
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println(err)
	}

	return &Meizam{
		connString: connString,
		db:         db,
	}
}

func (meizam *Meizam) GetUserState(userId int, facebookID string) (stateID int, lastMatchID int, lastDirection int) {
	query := fmt.Sprintf("exec spGetBotUserState %s, %d", facebookID, userId)
	res, err := meizam.db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Close()
	for res.Next() {
		var facebookID int64
		var uid int
		res.Scan(&facebookID, &uid, &stateID, &lastMatchID, &lastDirection)
		fmt.Println(facebookID, uid, stateID)
	}

	if err != nil {
		fmt.Println(err)
	}
	return stateID, lastMatchID, lastDirection
}

func (meizam *Meizam) UpdateUserState(userId int, stateID int, lastMatchID int, lastDirection int) error {
	query := fmt.Sprintf("exec spUpdateBotUserState %d, %d, %d, %d", userId, stateID, lastMatchID, lastDirection)
	res, err := meizam.db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Close()
	return err
}

func (meizam *Meizam) GetNextPredictionsToFill(userID int, tournamentID int, amount int) []*model.NextGame {
	query := fmt.Sprintf("exec spGetNextPredictionsToFill %d, %d, %d", userID, tournamentID, amount)
	res, err := meizam.db.Query(query)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Close()
	nextGames := []*model.NextGame{}
	for res.Next() {
		next := &model.NextGame{}
		res.Scan(&next.MatchID, &next.HomeTeam, &next.HomeTeamID, &next.AwayTeam, &next.AwayTeamID, &next.StartTime)
		nextGames = append(nextGames, next)
	}

	if err != nil {
		fmt.Println(err)
	}
	return nextGames
}

func (meizam *Meizam) SendDirectionPrediction(userID int, tournamentID int, matchID int, direction int) error {
	query := fmt.Sprintf("exec spAddSingleGameDirectionToTournament %d, %d, %d, %d", tournamentID, userID, matchID, direction)
	res, err := meizam.db.Query(query)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Close()
	return err
}

func (meizam *Meizam) SendScorePrediction(userID int, tournamentID int, matchID int, homeTeamScore int, awayTeamScore int) error {
	query := fmt.Sprintf("exec spAddSingleGameScoreToTournament %d, %d, %d, %d, %d", tournamentID, userID, matchID, homeTeamScore, awayTeamScore)
	res, err := meizam.db.Query(query)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Close()
	return err
}

func (meizam *Meizam) GetMatchDetails(matchID int) (homeTeamID int, awayTeamID int) {
	query := fmt.Sprintf("exec spGetMatchDetails %d", matchID)
	res, err := meizam.db.Query(query)
	if err != nil {
		return 0, 0
	}
	defer res.Close()
	for res.Next() {
		res.Scan(&homeTeamID, &awayTeamID)
		return homeTeamID, awayTeamID
	}
	return 0, 0
}

func (meizam *Meizam) GetUserId(facebookID string) (userID int) {
	query := fmt.Sprintf("exec spGetUserForFacebookBot %s", facebookID)
	res, err := meizam.db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Close()
	for res.Next() {
		res.Scan(&userID)
	}
	return userID
}
