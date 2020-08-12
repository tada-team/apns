package apns

import (
	"fmt"
)

type WebPush struct {
	// The title of the notification.
	Title string

	// The body of the notification.
	Body string

	// Optional. The label of the action button
	Action string

	// The url-args key specifies an array of values that are paired with the placeholders inside the urlFormatString
	// value of your website.json file. The url-args key must be included. The number of elements in the array must
	// match the number of placeholders in the urlFormatString value and the order of the placeholders in the URL
	// format string determines the order of the values supplied by the url-args array. The number of placeholders
	// may be zero, in which case the array should be empty. However, it is common practice to always include at
	// least one argument so that the user is directed to a web page specific to the notification received.
	UrlArgs []string

	Token string
}

func (p WebPush) Send(c *Config) (r Result) {
	url := fmt.Sprintf(urlMask, c.Host, p.Token)

	req := new(struct {
		Aps aps `json:"aps"`
	})

	req.Aps.UrlArgs = &p.UrlArgs
	req.Aps.Alert = &alert{
		Title:  p.Title,
		Body:   p.Body,
		Action: p.Action,
	}

	client, err := c.getSafariClient()
	if err != nil {
		r.Code = FailNow
		r.Error = err
		return
	}

	return c.Send(url, req, Headers{}, client)
}
