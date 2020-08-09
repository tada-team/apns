package firebase

type Preset struct {
	Name    string
	ApiKey  string
	Project string
}

func (p Preset) getAccessToken() string {
	return ""
}
