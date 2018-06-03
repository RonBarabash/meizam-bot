package model

type DirectionButton struct {
	title   string
	payload string
}

func (directionButton *DirectionButton) Title() string {
	return directionButton.title
}

func (directionButton *DirectionButton) Payload() string {
	return directionButton.payload
}

func (directionButton *DirectionButton) Type() string {
	return "postback"
}

func NewDirectionButton(title string, payload string) *DirectionButton {
	return &DirectionButton{
		title,
		payload,
	}
}
