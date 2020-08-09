package firebase

type Priority string

const (
	NormalPriority = Priority("normal")
	HighPriority   = Priority("high")
)

type NotificationPriority string

const (
	NotificationPriorityUnspecified = NotificationPriority("PRIORITY_UNSPECIFIED")

	// Lowest notification priority. Notifications with this PRIORITY_MIN might not be shown to the user
	// except under special circumstances, such as detailed notification logs.
	NotificationPriorityMin = NotificationPriority("PRIORITY_MIN")

	// Lower notification priority. The UI may choose to show the notifications smaller, or at a different
	// position in the list, compared with notifications with PRIORITY_DEFAULT.
	NotificationPriorityLow = NotificationPriority("PRIORITY_LOW")

	// Default notification priority. If the application does not prioritize its own notifications,
	// use this value for all notifications.
	NotificationPriorityDefault = NotificationPriority("PRIORITY_DEFAULT")

	// Higher notification priority. Use this for more important notifications or alerts. The UI may choose
	// to show these notifications larger, or at a different position in the notification lists, compared
	// with notifications with PRIORITY_DEFAULT.
	NotificationPriorityHigh = NotificationPriority("PRIORITY_HIGH")

	// Highest notification priority. Use this for the application's most important items that require
	// the user's prompt attention or input.
	NotificationPriorityMax = NotificationPriority("PRIORITY_MAX")
)

type NotificationVisibility string

const (
	// If unspecified, default to Visibility.PRIVATE.
	NotificationVisibilityUnspecified = NotificationVisibility("VISIBILITY_UNSPECIFIED")

	// Show this notification on all lockscreens, but conceal sensitive or private information on secure lockscreens.
	NotificationVisibilityPrivate = NotificationVisibility("PRIVATE")

	// Show this notification in its entirety on all lockscreens.
	NotificationVisibilityPublic = NotificationVisibility("PUBLIC")

	// Do not reveal any part of this notification on a secure lockscreen.
	NotificationVisibilitySecret = NotificationVisibility("SECRET")
)

type Color struct {
	Red   float32 `json:"red"`
	Green float32 `json:"green"`
	Blue  float32 `json:"blue"`
	Alpha float32 `json:"alpha"`
}

type LightSettings struct {
	Color            Color  `json:"color"`
	LightOnDuration  string `json:"light_on_duration"`
	LightOffDuration string `json:"light_off_duration"`
}

type AndroidNotification struct {
	BaseNotification

	// The key to the title string in the app's string resources to use to localize the title text to the
	// user's current localization.
	TitleLocKey string `json:"title_loc_key,omitempty"`

	// Variable string values to be used in place of the format specifiers in title_loc_key to use to
	// localize the title text to the user's current localization.
	TitleLocArgs []string `json:"title_loc_args,omitempty"`

	// The key to the body string in the app's string resources to use to localize the body text to the user's
	// current localization.
	BodyLocKey string `json:"body_loc_key,omitempty"`

	// Variable string values to be used in place of the format specifiers in body_loc_key to use to localize
	// the body text to the user's current localization.
	BodyLocArgs []string `json:"body_loc_args,omitempty"`

	// The notification's icon. Sets the notification icon to myicon for drawable resource myicon.
	// If you don't send this key in the request, FCM displays the launcher icon specified in your app manifest.
	Icon string `json:"icon,omitempty"`

	// The notification's icon color, expressed in #rrggbb format.
	Color string `json:"color,omitempty"`

	// The sound to play when the device receives the notification. Supports "default" or the filename
	// of a sound resource bundled in the app. Sound files must reside in /res/raw/.
	Sound string `json:"sound,omitempty"`

	// Identifier used to replace existing notifications in the notification drawer. If not specified, each
	// request creates a new notification. If specified and a notification with the same tag is already
	// being shown, the new notification replaces the existing one in the notification drawer.
	Tag string `json:"tag,omitempty"`

	// The action associated with a user click on the notification. If specified, an activity with a matching
	// intent filter is launched when a user clicks on the notification.
	ClickAction string `json:"click_action,omitempty"`

	// The notification's channel id (new in Android O). The app must create a channel with this channel ID
	// before any notification with this channel ID is received. If you don't send this channel ID in the
	// request, or if the channel ID provided has not yet been created by the app, FCM uses the
	// channel ID specified in the app manifest.
	ChannelId string `json:"channel_id,omitempty"`

	// Sets the "ticker" text, which is sent to accessibility services. Prior to API level 21 (Lollipop),
	// sets the text that is displayed in the status bar when the notification first arrives.
	Ticker string `json:"ticker,omitempty"`

	// When set to false or unset, the notification is automatically dismissed when the user clicks
	// it in the panel. When set to true, the notification persists even when the user clicks it.
	Sticky bool `json:"sticky,omitempty"`

	// Set whether or not this notification is relevant only to the current device. Some notifications
	// can be bridged to other devices for remote display, such as a Wear OS watch. This hint can be
	// set to recommend this notification not be bridged.
	LocalOnly bool `json:"local_only,omitempty"`

	// Set the relative priority for this notification. Priority is an indication of how much of the user's
	// attention should be consumed by this notification. Low-priority notifications may be hidden from the user
	// in certain situations, while the user might be interrupted for a higher-priority notification.
	// The effect of setting the same priorities may differ slightly on different platforms. Note
	// this priority differs from AndroidMessagePriority. This priority is processed by the client
	// after the message has been delivered, whereas AndroidMessagePriority is an FCM concept that controls
	// when the message is delivered.
	NotificationPriority NotificationPriority `json:"notification_priority,omitempty"`

	// If set to true, use the Android framework's default sound for the notification.
	// Default values are specified in config.xml.
	DefaultSound bool `json:"default_sound,omitempty"`

	// If set to true, use the Android framework's default vibrate pattern for the notification. Default values
	// are specified in config.xml. If default_vibrate_timings is set to true and vibrate_timings is also set,
	// the default value is used instead of the user-specified vibrate_timings.
	DefaultVibrateTimings bool `json:"default_vibrate_timings,omitempty"`

	// If set to true, use the Android framework's default LED light settings for the notification. Default
	// values are specified in config.xml. If default_light_settings is set to true and light_settings is also set,
	// the user-specified light_settings is used instead of the default value.
	DefaultLightSettings bool `json:"default_light_settings,omitempty"`

	// Set the vibration pattern to use. Pass in an array of protobuf.Duration to turn on or off the vibrator.
	// The first value indicates the Duration to wait before turning the vibrator on. The next value indicates
	// the Duration to keep the vibrator on. Subsequent values alternate between Duration to turn the vibrator
	// off and to turn the vibrator on. If vibrate_timings is set and default_vibrate_timings is set to true,
	// the default value is used instead of the user-specified vibrate_timings.
	VibrateTimings []string `json:"vibrate_timings,omitempty"`

	// Sphere of visibility of this notification, which affects how and when the SystemUI reveals the
	// notification's presence and contents in untrusted situations (namely, on the secure lockscreen).
	//The default level, VISIBILITY_PRIVATE, behaves exactly as notifications have always done on Android:
	// The notification's icon and tickerText (if available) are shown in all situations, but the contents
	// are only available if the device is unlocked for the appropriate user. A more permissive policy can
	// be expressed by VISIBILITY_PUBLIC; such a notification can be read even in an "insecure" context
	// (that is, above a secure lockscreen). To modify the public version of this notification for example,
	// to redact some portions?see Builder#setPublicVersion(Notification). Finally, a notification can be
	// made VISIBILITY_SECRET, which will suppress its icon and ticker until the user has bypassed the lockscreen.
	Visibility NotificationVisibility `json:"visibility,omitempty"`

	// Sets the number of items this notification represents. May be displayed as a badge count for launchers
	// that support badging. For example, this might be useful if you're using just one notification to
	// represent multiple new messages but you want the count here to represent the number
	// of total new messages. If zero or unspecified, systems that support badging use the default,
	// which is to increment a number displayed on the long-press menu each time a new notification arrives.
	NotificationCount int `json:"notification_count,omitempty"`

	// Settings to control the notification's LED blinking rate and color if LED is available on the device.
	// The total blinking time is controlled by the OS.
	LightSettings *LightSettings `json:"light_settings,omitempty"`
}

type AndroidConfig struct {
	// An identifier of a group of messages that can be collapsed, so that only the last message gets
	//sent when delivery can be resumed. A maximum of 4 different collapse keys is allowed at any given time.
	CollapseKey string `json:"collapse_key,omitempty"`

	// Message priority. Can take "normal" and "high" values.
	Priority Priority `json:"priority,omitempty"`

	// How long (in seconds) the message should be kept in FCM storage if the device is offline.
	// The maximum time to live supported is 4 weeks, and the default value is 4 weeks if not set.
	// Set it to 0 if want to send the message immediately. In JSON format, the Duration type
	// is encoded as a string rather than an object, where the string ends in the suffix "s"
	// (indicating seconds) and is preceded by the number of seconds, with nanoseconds
	// expressed as fractional seconds. For example, 3 seconds with 0 nanoseconds should
	// be encoded in JSON format as "3s", while 3 seconds and 1 nanosecond should be expressed
	// in JSON format as "3.000000001s". The ttl will be rounded down to the nearest second.
	TTL int64 `json:"ttl,omitempty"`

	// Package name of the application where the registration token must match in order to receive the message.
	RestrictedPackageName string `json:"restricted_package_name"`

	//Arbitrary key/value payload. If present, it will override google.firebase.fcm.v1.Message.data.
	Data map[string]string `json:"data,omitempty"`

	// Notification to send to android devices.
	Notification *AndroidNotification `json:"notification,omitempty"`

	// Options for features provided by the FCM SDK for Android.
	FcmOptions *FcmOptions `json:"fcm_options,omitempty"`

	// If set to true, messages will be allowed to be delivered to the app while the device is in direct boot mode.
	// https://developer.android.com/training/articles/direct-boot
	DirectBootOk bool `json:"direct_boot_ok,omitempty"`
}
