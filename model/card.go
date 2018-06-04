package model

import "github.com/RonBarabash/meizam-bot/interfaces"

type Card struct {
	title    string
	subTitle string
	imageURL string
	buttons  []messaging.IButton
}

func (gameCard *Card) Title() string {
	return gameCard.title
}

func (gameCard *Card) Subtitle() string {
	return gameCard.subTitle
}

func (gameCard *Card) ImageURL() string {
	return gameCard.imageURL
}

func (gameCard *Card) Buttons() []messaging.IButton {
	return gameCard.buttons
}

func NewCard(title string, subTitle string, imageURL string, buttons []messaging.IButton) *Card {
	return &Card{
		title,
		subTitle,
		imageURL,
		buttons,
	}
}
