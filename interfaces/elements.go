package messaging

type IButton interface{
	Type() string
	Title() string
	Payload() string
}

type ICard interface {
	Title() string
	Subtitle() string
	ImageURL() string
	Buttons() []IButton
}