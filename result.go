package apns

type Result struct {
	Code          ResultCode
	Error         error
	DebugRequest  string
	DebugResponse string
}

type ResultCode int

const (
	Ok ResultCode = iota
	FailNow
	RetryNow
	RetryLater
	InvalidConfig
)
