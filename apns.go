package firebase

type ApnsFcmOptions struct {
	AnalyticsLabel string `json:"analytics_label"`
	Image          string `json:"image"`
}

type ApnsConfig struct {
	// HTTP request headers defined in Apple Push Notification Service. Refer to APNs request headers
	// for supported headers, e.g. "apns-priority": "10".
	Headers map[string]string `json:"headers,omitempty"`

	// APNs payload as a JSON object, including both aps dictionary and custom payload. See Payload Key Reference.
	Payload map[string]string `json:"payload,omitempty"`

	// Options for features provided by the FCM SDK for Web.
	FcmOptions *ApnsFcmOptions `json:"fcm_options"`
}
