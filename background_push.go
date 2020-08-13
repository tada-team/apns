package apns

import (
	"fmt"
)

type BackgroundPush struct {
	Category string
	ThreadId string
	Badge    *int
	Data     map[string]interface{}
	Token    string
}

func (p BackgroundPush) Send(c *Config, h *Headers) (r Result) {
	url := fmt.Sprintf(urlMask, c.Host, p.Token)

	req := make(map[string]interface{})
	for k, v := range p.Data {
		req[k] = v
	}

	req["apns"] = aps{
		ContentAvailable: 1,
		Badge:            p.Badge,
		Category:         p.Category,
		ThreadId:         p.ThreadId,
	}

	if h == nil {
		h = new(Headers)
	}
	h.PushType = PushTypeBackground

	return c.Send(url, req, *h, nil)
}
