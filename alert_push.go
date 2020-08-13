package apns

import (
	"fmt"
)

type AlertPush struct {
	Title    string
	Subtitle string
	Body     string
	Badge    *int
	Data     map[string]interface{}
	Token    string
}

func (p AlertPush) Send(c *Config, h *Headers) (r Result) {
	url := fmt.Sprintf(urlMask, c.Host, p.Token)

	req := make(map[string]interface{})
	for k, v := range p.Data {
		req[k] = v
	}

	req["apns"] = aps{
		Badge:    p.Badge,
		Sound:    "default",
		Category: "QuickReply",
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
