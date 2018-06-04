package model

import "github.com/RonBarabash/meizam-bot/interfaces"

type GameCard struct {
	title    string
	subTitle string
	imageURL string
	buttons  []messaging.IButton
}

func (gameCard *GameCard) Title() string {
	return gameCard.title
}

func (gameCard *GameCard) Subtitle() string {
	return gameCard.subTitle
}

func (gameCard *GameCard) ImageURL() string {
	return gameCard.imageURL
}

func (gameCard *GameCard) Buttons() []messaging.IButton {
	return gameCard.buttons
}

func NewGameCard(title string, subTitle string, imageURL string, buttons []messaging.IButton) *GameCard {
	return &GameCard{
		title,
		subTitle,
		imageURL,
		buttons,
	}
}
