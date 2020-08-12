package apns

type alert struct {
	// The title of the notification. Apple Watch displays this string in the short look notification interface.
	// Specify a string that is quickly understood by the user.
	Title string `json:"title,omitempty"`

	// Additional information that explains the purpose of the notification.
	Subitle string `json:"subtitle,omitempty"`

	// The content of the alert message.
	Body string `json:"body"`

	// Safari only.
	Action string `json:"action,omitempty"`

	// The name of the launch image file to display. If the user chooses to launch your app, the contents
	// of the specified image or storyboard file are displayed instead of your app's normal launch image.
	LaunchImage string `json:"launch-image,omitempty"`

	// The key for a localized title string. Specify this key instead of the title key to retrieve the title
	// from your app’s Localizable.strings files. The value must contain the name of a key in your strings file.
	TitleLocKey string `json:"title-loc-key,omitempty"`

	// An array of strings containing replacement values for variables in your title string. Each %@ character
	// in the string specified by the title-loc-key is replaced by a value from this array. The first item in
	// the array replaces the first instance of the %@ character in the string, the second item replaces the
	// second instance, and so on.
	TitleLocArgs *[]string `json:"title-loc-args,omitempty"`

	// The key for a localized subtitle string. Use this key, instead of the subtitle key, to retrieve the subtitle
	// from your app's Localizable.strings file. The value must contain the name of a key in your strings file.
	SubtitleLocKey string `json:"subtitle-loc-key,omitempty"`

	// An array of strings containing replacement values for variables in your title string. Each %@ character
	// in the string specified by subtitle-loc-key is replaced by a value from this array. The first item in
	// the array replaces the first instance of the %@ character in the string, the second item replaces
	// the second instance, and so on.
	SubtitleLocArgs *[]string `json:"subtitle-loc-args,omitempty"`

	// If a string is specified, the system displays an alert that includes the Close and View buttons.
	// The string is used as a key to get a localized string in the current localization to use for the
	// right button’s title instead of “View”.
	ActionLocKey string `json:"action-loc-key,omitempty"`

	// The key for a localized message string. Use this key, instead of the body key, to retrieve the message text
	// from your app's Localizable.strings file. The value must contain the name of a key in your strings file.
	LocKey string `json:"loc-key,omitempty"`

	// An array of strings containing replacement values for variables in your message text. Each %@ character
	// in the string specified by loc-key is replaced by a value from this array. The first item in the array
	// replaces the first instance of the %@ character in the string, the second item replaces the second
	// instance, and so on.
	LocArgs *[]string `json:"loc-args,omitempty"`
}

type aps struct {
	// The information for displaying an alert
	Alert *alert `json:"alert,omitempty"`

	// The number to display in a badge on your app’s icon. Specify 0 to remove the current badge, if any.
	Badge *int `json:"badge,omitempty"`

	// The name of a sound file in your app’s main bundle or in the Library/Sounds folder of your app’s
	// container directory. Specify the string "default" to play the system sound. Use this key for
	// regular notifications.
	Sound string `json:"sound,omitempty"`

	// An app-specific identifier for grouping related notifications. This value corresponds to the
	// threadIdentifier property in the UNNotificationContent object.
	ThreadId string `json:"thread-id,omitempty"`

	// The notification’s type. This string must correspond to the identifier of one of the UNNotificationCategory
	// objects you register at launch time.
	Category string `json:"category,omitempty"`

	// The background notification flag. To perform a silent background update, specify the value 1 and don't
	// include the alert, badge, or sound keys in your payload.
	ContentAvailable int `json:"content-available,omitempty"`

	// The notification service app extension flag. If the value is 1, the system passes the notification
	// to your notification service app extension before delivery. Use your extension to modify the
	// notification’s content.
	MutableContent int `json:"mutable-content,omitempty"`

	// The identifier of the window brought forward. The value of this key will be populated on the
	// UNNotificationContent object created from the push payload. Access the value using the UNNotificationContent
	// object's targetContentIdentifier property.
	TargetContentId string `json:"target-content-id,omitempty"`

	// Safari only
	UrlArgs *[]string `json:"url-args,omitempty"`
}
