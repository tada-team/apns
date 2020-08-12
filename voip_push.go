package apns

import (
	"fmt"
)

type Push interface {
	Send(c *Config)
}

type VoipPush struct {
	Title string
	Body  string
	Token string
}

func (p VoipPush) Send(c *Config, h *Headers) (r Result) {
	url := fmt.Sprintf(urlMask, c.Host, p.Token)

	req := new(struct {
		Aps aps `json:"aps"`
	})

	req.Aps.Alert = &alert{
		Title: p.Title,
		Body:  p.Body,
	}

	if h == nil {
		h = new(Headers)
	}
	h.PushType = PushTypeVoip

	return c.Send(url, req, *h, nil)
}
