package apns

import (
	"fmt"
)

type AlertPush struct {
	BackgroundPush
	Title    string
	Subtitle string
	Body     string
	Sound    string
}

func (p AlertPush) Send(c *Config, h *Headers) (r Result) {
	url := fmt.Sprintf(urlMask, c.Host, p.Token)

	req := make(map[string]interface{})
	for k, v := range p.Data {
		req[k] = v
	}

	if p.Sound == "" {
		p.Sound = "default"
	}

	req["apns"] = aps{
		Badge:    p.Badge,
		Category: p.Category,
		ThreadId: p.ThreadId,
		Sound:    p.Sound,
		Alert: &alert{
			Title:   p.Title,
			Subitle: p.Subtitle,
			Body:    p.Body,
		},
	}

	if h == nil {
		h = new(Headers)
	}
	h.PushType = PushTypeAlert

	return c.Send(url, req, *h, nil)
}
