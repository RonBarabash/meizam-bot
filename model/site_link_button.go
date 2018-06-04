package model

type SiteLinkButton struct {}

func (siteLinkButton *SiteLinkButton) Title() string {
	return "הרשם"
}

func (siteLinkButton *SiteLinkButton) Payload() string {
	return "http://www.meizam.club/Home/SingleGames"
}

func (siteLinkButton *SiteLinkButton) Type() string {
	return "web_url"
}

func NewSiteLinkButton() *SiteLinkButton {
	return &SiteLinkButton{}
}
