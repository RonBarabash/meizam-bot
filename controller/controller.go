package controller

import (
	"github.com/maciekmm/messenger-platform-go-sdk"
	"fmt"
	"github.com/meizam-bot/meizam"
	"github.com/meizam-bot/providers"
)

type Controller struct {
	meizam            *meizam.Meizam
	messengerProvider *providers.FacebookMessengerProvider
}

func (controller *Controller) BindMessageReceived() messenger.MessageReceivedHandler {
	return func(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
		fmt.Println("got message: " + msg.Text)
		userId := 56
		facebookID := opts.Sender.ID
		switch userState := controller.meizam.GetUserState(userId, facebookID); userState {
		case 1:
			//explain who u are
			controller.messengerProvider.SendSimpleMessage(facebookID, "אהלן! אני שימי ואני כאן לעזור לך לסיים במקום הראשון במיזם! 🤑  ")
			//send next games
			controller.meizam.GetNextPredictionsToFill(userId, 4, 3)
			//update to next state
			controller.meizam.UpdateUserState(userId, 1)

		default:
			fmt.Printf("Default")
		}

		//if err != nil {
		//	fmt.Println("error handling message: " + err.Error())
		//}
		fmt.Println("handled message: " + msg.Text)
	}
}

func NewController(meizam *meizam.Meizam, provider *providers.FacebookMessengerProvider) *Controller {
	return &Controller{meizam, provider}
}
