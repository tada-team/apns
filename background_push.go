package apns

import (
	"fmt"
)

type BackgroundPush struct {
	Badge *int
	Data  map[string]interface{}
	Token string
}

func (p BackgroundPush) Send(c *Config, h *Headers) (r Result) {
	url := fmt.Sprintf(urlMask, c.Host, p.Token)

	req := p.Data
	req["apns"] = aps{
		ContentAvailable: 1,
		Badge:            p.Badge,
	}

	if h == nil {
		h = new(Headers)
	}
	h.PushType = PushTypeBackground

	return c.Send(url, req, *h, nil)
}