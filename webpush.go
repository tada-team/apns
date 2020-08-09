package firebase

type WebpushFcmOptions struct {
	Link           string `json:"link"`
	AnalyticsLabel string `json:"analytics_label"`
}

type WebpushConfig struct {
	// HTTP headers defined in webpush protocol. Refer to Webpush protocol for supported headers, e.g. "TTL": "15"
	Headers map[string]string `json:"headers,omitempty"`

	// An object containing a list of "key": value pairs. Example: { "name": "wrench", "mass": "1.3kg", "count": "3" }
	Data map[string]string `json:"data,omitempty"`

	// Web Notification options as a JSON object.
	Notification interface{} `json:"notification,omitempty"`

	// Options for features provided by the FCM SDK for Web.
	FcmOptions *WebpushFcmOptions `json:"fcm_options"`
}
