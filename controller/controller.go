package controller

import (
	"github.com/maciekmm/messenger-platform-go-sdk"
	"fmt"
	"github.com/RonBarabash/meizam-bot/meizam"
	"github.com/RonBarabash/meizam-bot/providers"
	"github.com/RonBarabash/meizam-bot/model"
	"github.com/RonBarabash/meizam-bot/interfaces"
	"strings"
	"strconv"
)

type Controller struct {
	meizam            *meizam.Meizam
	messengerProvider *providers.FacebookMessengerProvider
}

func (controller *Controller) BindMessageReceived() messenger.MessageReceivedHandler {
	return func(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
		fmt.Println("got message: " + msg.Text)
		facebookID := opts.Sender.ID
		userId := controller.meizam.GetUserId(facebookID)
		if userId == 0 {
			stateID, _, _ := controller.meizam.GetUserState(userId, facebookID)
			if stateID == 1 {
				err := controller.messengerProvider.SendSimpleMessage(facebookID, "צריך להרשם תחילה דרך האתר  - כי אין לי מושג מי את/ה")
				if err != nil {
					fmt.Println(err)
				}
				buttons := []messaging.IButton{model.NewSiteLinkButton()}
				cards := []messaging.ICard{model.NewCard("הרשם לבוט באתר", "ואז תחזור אלי :)", "", buttons)}
				err = controller.messengerProvider.SendGenericTemplate(facebookID, map[string]string{}, cards)
				if err != nil {
					fmt.Println(err)
				}
			} else if stateID == 3 {
				controller.messengerProvider.SendSimpleMessage(facebookID, "עד שתרשם דרך האתר - אני לא אזכור דבר ממה שאתה אומר :)")
			}
		}
		userState, lastMatchID, lastDirection := controller.meizam.GetUserState(userId, facebookID)
		switch userState {
		case 1:
			//explain who u are
			controller.messengerProvider.SendSimpleMessage(facebookID, "אהלן! אני ויטליק ואני כאן לעזור לך לסיים במקום הראשון במיזם! 🤑  ")
			//send next games
			controller.sendGames(userId, facebookID)
			//update to next state
			controller.meizam.UpdateUserState(userId, 2, 0, 0)
		case 2:
			controller.messengerProvider.SendSimpleMessage(facebookID, "בגדול אני מחכה שתבחר כיוון...")
		case 3:
			homeTeamID, _ := controller.meizam.GetMatchDetails(lastMatchID)
			parts := strings.Split(strings.TrimSpace(msg.Text), "-")
			if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
				controller.messengerProvider.SendSimpleMessage(facebookID, "נא לנסות שוב")
				fmt.Println("dont understand this: " + msg.Text)
				return
			}
			firstScore, errFirstPart := strconv.Atoi(parts[0])
			secondScore, errSecondPart := strconv.Atoi(parts[1])
			if (errFirstPart != nil || errSecondPart != nil) {
				controller.messengerProvider.SendSimpleMessage(facebookID, "נא לענות בפורמט שביקשתי בלבד, ורק במספרים. לדוגמה: 2-1")
				fmt.Println("dont understand this: " + msg.Text)
				return
			}
			if lastDirection == 0 {
				if (firstScore == secondScore) {
					controller.meizam.SendScorePrediction(userId, 4, lastMatchID, firstScore, firstScore)
				} else {
					controller.messengerProvider.SendSimpleMessage(facebookID, "באת לשגע אותי? לא הבנתי - אז לא שמרתי את התוצאה")
				}
			} else {
				if lastDirection == homeTeamID {
					if firstScore > secondScore {
						controller.meizam.SendScorePrediction(userId, 4, lastMatchID, firstScore, secondScore)
					} else {
						controller.meizam.SendScorePrediction(userId, 4, lastMatchID, secondScore, firstScore)
					}
				} else {
					if firstScore > secondScore {
						controller.meizam.SendScorePrediction(userId, 4, lastMatchID, secondScore, firstScore)
					} else {
						controller.meizam.SendScorePrediction(userId, 4, lastMatchID, firstScore, secondScore)
					}
				}
				controller.meizam.UpdateUserState(userId, 2, 0, 0)
				controller.sendGames(userId, facebookID)
			}

		default:
			fmt.Printf("Default")
		}
		fmt.Println("handled message: " + msg.Text)
	}
}

func (controller *Controller) sendGames(userId int, facebookID string) {
	games := controller.meizam.GetNextPredictionsToFill(userId, 4, 3)
	gameCards := []messaging.ICard{}
	for _, game := range games {
		buttons := []messaging.IButton{}

		buttons = append(buttons, model.NewDirectionButton(game.HomeTeam, fmt.Sprintf("direction-%d-%d", game.MatchID, game.HomeTeamID)))
		buttons = append(buttons, model.NewDirectionButton(game.AwayTeam, fmt.Sprintf("direction-%d-%d", game.MatchID, game.AwayTeamID)))
		buttons = append(buttons, model.NewDirectionButton("תיקו", fmt.Sprintf("direction-%d-%d", game.MatchID, 0)))

		gameCard := model.NewCard(fmt.Sprintf("%s-%s", game.HomeTeam, game.AwayTeam), "איך יסתיים?", "", buttons)
		gameCards = append(gameCards, gameCard)
	}
	quickReplies := make(map[string]string)
	controller.messengerProvider.SendGenericTemplate(facebookID, quickReplies, gameCards)
}

func (controller *Controller) BindPostbackReceived() messenger.PostbackHandler {
	return func(event messenger.Event, opts messenger.MessageOpts, postback messenger.Postback) {
		facebookID := opts.Sender.ID
		fakeMid := fmt.Sprintf("postback_%d", event.ID)
		fmt.Println("got postback: " + fakeMid)
		parts := strings.Split(postback.Payload, "-")
		userId := controller.meizam.GetUserId(facebookID)
		switch parts[0] {
		case "direction":
			matchID, _ := strconv.Atoi(parts[1])
			direction, _ := strconv.Atoi(parts[2])
			err := controller.meizam.SendDirectionPrediction(userId, 4, matchID, direction)
			if err != nil {
				fmt.Println(err)
			} else {
				controller.meizam.UpdateUserState(userId, 3, matchID, direction)
				controller.messengerProvider.SendSimpleMessage(facebookID, fmt.Sprintf("ומה תהיה התוצאה? "))
				controller.messengerProvider.SendSimpleMessage(facebookID, fmt.Sprintf("אני אבין אם יכתבו לי משהו כמו 3-0"))
			}
		}
	}
}

func (controller *Controller) BindAuthentication() messenger.AuthenticationHandler {
	return func(event messenger.Event, opts messenger.MessageOpts, optin *messenger.Optin) {
		fmt.Println(fmt.Sprintf("got optin for u_id: %v", optin.Ref))
		userId, _ := strconv.Atoi(optin.Ref)
		facebookID := opts.Sender.ID
		userState, _, _ := controller.meizam.GetUserState(userId, facebookID)
		if userState == 1 {
			controller.messengerProvider.SendSimpleMessage(facebookID, "אהלן! אני ויטליק ואני כאן לעזור לך לסיים במקום הראשון במיזם! 🤑  ")
		}
		//send next games
		controller.sendGames(userId, facebookID)
		//update to next state
		controller.meizam.UpdateUserState(userId, 2, 0, 0)
	}
}
func NewController(meizam *meizam.Meizam, provider *providers.FacebookMessengerProvider) *Controller {
	return &Controller{meizam, provider}
}
