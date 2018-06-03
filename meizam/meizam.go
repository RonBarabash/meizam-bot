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

func (meizam *Meizam) GetUserState(userId int, facebookID string) (stateID int) {
	query := fmt.Sprintf("exec spGetBotUserState %s, %d", facebookID, userId)
	res, err := meizam.db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Close()
	for res.Next() {
		var facebookID int64
		var uid int
		res.Scan(&facebookID, &uid, &stateID)
		fmt.Println(facebookID, uid, stateID)
	}

	if err != nil {
		fmt.Println(err)
	}
	return stateID
}

func (meizam *Meizam) UpdateUserState(userId int, stateID int) error {
	query := fmt.Sprintf("exec spUpdateBotUserState %d, %d", userId, stateID)
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
