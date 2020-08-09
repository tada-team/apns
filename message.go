package firebase

type BaseNotification struct {
	// The notification's title
	Title string `json:"title,omitempty"`

	// The notification's body text
	Body string `json:"body,omitempty"`

	// Contains the URL of an image that is going to be downloaded on the device and displayed in a notification.
	// JPEG, PNG, BMP have full support across platforms. Animated GIF and video only work on iOS. WebP and HEIF
	// have varying levels of support across platforms and platform versions. Android has 1MB image size limit.
	Image string `json:"image"`
}

type Message struct {
	// Output Only. The identifier of the message sent, in the format of projects/*/messages/{message_id}.
	Name string `json:"name,omitempty"`

	// Arbitrary key/value payload. If present, it will override google.firebase.fcm.v1.Message.data.
	Data map[string]string `json:"data,omitempty"`

	// Input only. Basic notification template to use across all platforms.
	Notification *BaseNotification `json:"notification,omitempty"`

	// Input only. Android specific options for messages sent through FCM connection server.
	Android *AndroidConfig `json:"android,omitempty"`

	// Input only. Webpush protocol options.
	Webpush *WebpushConfig `json:"webpush,omitempty"`

	// Input only. Apple Push Notification Service specific options.
	Apns *ApnsConfig `json:"apns,omitempty"`

	// Input only. Template for FCM SDK feature options to use across all platforms.
	FcmOptions *FcmOptions `json:"fcm_options,omitempty"`

	// Union field target can be only one of the following:
	Token     string `json:"token,omitempty"`     // Registration token to send a message to.
	Topic     string `json:"topic,omitempty"`     // Topic name to send a message to, e.g. "weather". Note: "/topics/" prefix should not be provided.
	Condition string `json:"condition,omitempty"` // Condition to send a message to, e.g. "'foo' in topics && 'bar' in topics".
	// End of list of possible types for union field target.
}
