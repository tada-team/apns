package apns

import (
	"fmt"
)

type AlertPush struct {
	Title    string
	Subtitle string
	Body     string
	Badge    *int
	Token    string
}

func (p AlertPush) Send(c *Config, h *Headers) (r Result) {
	url := fmt.Sprintf(urlMask, c.Host, p.Token)

	req := new(struct {
		Aps aps `json:"aps"`
	})
	req.Aps.Sound = "default"
	req.Aps.Category = "QuickReply"
	req.Aps.Badge = p.Badge
	req.Aps.Alert = &alert{
		Title:   p.Title,
		Subitle: p.Subtitle,
		Body:    p.Body,
	}

	if h == nil {
		h = new(Headers)
	}
	h.PushType = PushTypeAlert

	return c.Send(url, req, *h, nil)
}
