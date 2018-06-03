package controller

import (
	"github.com/maciekmm/messenger-platform-go-sdk"
	"fmt"
)

type Controller struct {}


func (controller *Controller) BindMessageReceived() messenger.MessageReceivedHandler {
	return func(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
		fmt.Println("got message: " + msg.Text)
		//if err != nil {
		//	fmt.Println("error handling message: " + err.Error())
		//}
		fmt.Println("handled message: " + msg.Text)
	}
}


func NewController() *Controller {
	return &Controller{}
}