package firebase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type resultCode int

type Result struct {
	Code          resultCode
	Error         error
	DebugRequest  string
	DebugResponse string
}

const (
	Ok resultCode = iota + 1
	Fatal
	Retry
	InvalidPreset
)

type request struct {
	Message      Message `json:"message"`
	ValidateOnly bool    `json:"validate_only,omitempty"`
}

type response struct {
	Message
	ErrorCode string `json:"error_code,omitempty"`
}

func (p Preset) Send(m Message) (r Result) {
	reqJson, err := json.Marshal(request{Message: m})
	if err != nil {
		r.Error = errors.Wrap(err, "json marshal fail")
		r.Code = Fatal
		return
	}

	url := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", p.Project)
	reqHttp, err := http.NewRequest("POST", url, bytes.NewReader(reqJson))
	r.DebugRequest = fmt.Sprintf("url: %s\nbody: %s", url, string(reqJson))
	if err != nil {
		r.Error = err
		r.Code = Retry
		return
	}

	reqHttp.Header.Set("Content-Type", "application/json")
	reqHttp.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.getAccessToken()))

	rawResp, err := http.DefaultClient.Do(reqHttp)
	if err != nil {
		r.Error = err
		r.Code = Retry
		return
	}
	defer rawResp.Body.Close()

	body, err := ioutil.ReadAll(rawResp.Body)
	if err != nil {
		r.Error = err
		r.Code = Retry
		return
	}

	r.DebugResponse = fmt.Sprintf("code: %d\nbody: %s", rawResp.StatusCode, string(body))

	resp := response{}
	if err := json.Unmarshal(body, &resp); err != nil {
		r.Error = errors.Wrap(err, "json error")
		r.Code = Fatal
		return
	}

	switch rawResp.StatusCode {
	case 200:
		r.Code = Ok
	case 400: // Request parameters were invalid
		r.Error = fmt.Errorf("firebase: invalid params: %s", body)
		r.Code = Fatal
	case 401:
		r.Error = fmt.Errorf("firebase: APNs certificate or web push auth key was invalid or missing")
		r.Code = InvalidPreset
	case 403:
		r.Error = fmt.Errorf("firebase: authenticated sender ID is different from the sender ID for the registration token")
		r.Code = InvalidPreset
	case 404:
		r.Error = fmt.Errorf("firebase: app instance was unregistered from FCM")
		r.Code = InvalidPreset
	case 429:
		r.Error = fmt.Errorf("firebase: sending limit exceeded for the message target")
		r.Code = Fatal
	case 500:
		r.Error = fmt.Errorf("firebase: unknown internal error occurred")
		r.Code = Retry
	case 503:
		r.Error = fmt.Errorf("firebase: server is overloaded")
		r.Code = Retry
	default:
		r.Error = fmt.Errorf("firebase: unknow code code %d body: %s", rawResp.StatusCode, body)
		r.Code = Fatal
	}
	return
}
