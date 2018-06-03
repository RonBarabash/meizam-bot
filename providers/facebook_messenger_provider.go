package providers

import (
	"github.com/maciekmm/messenger-platform-go-sdk"
	"github.com/meizam-bot/interfaces"
	"github.com/maciekmm/messenger-platform-go-sdk/template"
	"errors"
)

type FacebookMessengerProvider struct {
	messengerClient *messenger.Messenger
}

func NewFacebookMessengerProvider(messengerClient *messenger.Messenger) *FacebookMessengerProvider {
	return &FacebookMessengerProvider{messengerClient: messengerClient}
}

func (messengerProvider *FacebookMessengerProvider) SendSimpleMessage(externalUserID, text string) (err error) {
	_, err = messengerProvider.messengerClient.SendSimpleMessage(externalUserID, text)
	return err
}

//func (messengerProvider *FacebookMessengerProvider) SendQuickReplyMessage(externalUserID, text string, quickReplies map[string]string) (err error) {
//	recipient := messenger.Recipient{ID: externalUserID}
//	message := messenger.SendMessage{Text: text}
//	messageQuery := messenger.MessageQuery{Recipient: recipient, Message: message}
//	for title, payload := range quickReplies {
//		err = messageQuery.QuickReply(title, payload)
//		if err != nil {
//			return err
//		}
//	}
//	_, err = messengerProvider.messengerClient.SendMessage(messageQuery)
//	return err
//}
//
//func (messengerProvider *FacebookMessengerProvider) SendImage(externalUserID, imageURL string) (err error) {
//	recipient := messenger.Recipient{ID: externalUserID}
//	message := messenger.SendMessage{}
//	messageQuery := messenger.MessageQuery{Recipient: recipient, Message: message}
//
//	messageQuery.Image(imageURL)
//	_, err = messengerProvider.messengerClient.SendMessage(messageQuery)
//	return err
//}

func (messengerProvider *FacebookMessengerProvider) SendGenericTemplate(externalUserID string, quickReplies map[string]string, cards []messaging.ICard) error {
	messageQuery := messenger.MessageQuery{}
	messageQuery.RecipientID(externalUserID)

	for _, card := range cards {
		genericTemplate, err := castCardToGenericTemplate(card)
		if err != nil {
			return err
		}
		err = messageQuery.Template(genericTemplate)
		if err != nil {
			return err
		}
	}
	//for title, payload := range quickReplies {
	//	if err := messageQuery.QuickReply(title, payload); err != nil {
	//		return err
	//	}
	//}
	_, err := messengerProvider.messengerClient.SendMessage(messageQuery)
	return err
}

func castButtonsToFacebookButtons(buttons []messaging.IButton) ([]template.Button, error) {
	facebookButtons := []template.Button{}
	for _, button := range buttons {
		facebookButtonType, err := castButtonType(button.Type())
		if err != nil {
			return nil, err
		}
		facebookButton := template.Button{Type: facebookButtonType, Title: button.Title(), Payload: button.Payload()}
		facebookButtons = append(facebookButtons, facebookButton)
	}
	return facebookButtons, nil
}

func castCardToGenericTemplate(card messaging.ICard) (template.GenericTemplate, error) {
	buttons, err := castButtonsToFacebookButtons(card.Buttons())
	if err != nil {
		return template.GenericTemplate{}, err
	}
	genericTemplate := template.GenericTemplate{Title: card.Title(),
		Subtitle: card.Subtitle(),
		ImageURL: card.ImageURL(),
		Buttons: buttons,
	}
	return genericTemplate, nil
}

func castButtonType(buttonType string) (template.ButtonType, error) {
	switch buttonType {
	case "web_url":
		return template.ButtonTypeWebURL, nil
	case "postback":
		return template.ButtonTypePostback, nil
	case "phone_number":
		return template.ButtonTypePhoneNumber, nil
	case "account_link":
		return template.ButtonTypeAccountLink, nil
	case "account_unlink":
		return template.ButtonTypeAccountUnlink, nil
	default:
		return "", errors.New("unkown button type: " + buttonType)
	}
}
