package apns

import (
	"fmt"
	"time"
)

type PushType string

const (
	// Use the alert push type for notifications that trigger a user interaction—for example, an alert, badge,
	// or sound. If you set this push type, the apns-topic header field must use your app’s bundle ID as the topic.
	PushTypeAlert = PushType("alert")

	// Use the background push type for notifications that deliver content in the background, and don’t trigger any
	// user interactions. If you set this push type, the apns-topic header field must use your app’s bundle ID
	// as the topic. For more information, see Pushing Background Updates to Your App.
	PushTypeBackground = PushType("background")

	// Use the voip push type for notifications that provide information about an incoming Voice-over-IP (VoIP) call.
	// For more information, see Responding to VoIP Notifications from PushKit.
	// If you set this push type, the apns-topic header field must use your app’s bundle ID with .voip appended
	// to the end. If you’re using certificate-based authentication, you must also register the certificate
	// for VoIP services. The topic is then part of the 1.2.840.113635.100.6.3.4 or 1.2.840.113635.100.6.3.6 extension.
	PushTypeVoip = PushType("voip")

	// Use the complication push type for notifications that contain update information for a watchOS app’s
	// complications. If you set this push type, the apns-topic header field must use your app’s bundle ID  with
	// .complication appended to the end. If you’re using certificate-based authentication, you must also register
	// the certificate for WatchKit services. The topic is then part of the 1.2.840.113635.100.6.3.6 extension.
	PushTypeComplication = PushType("complication")

	// Use the fileprovider push type to signal changes to a File Provider extension. If you set this push type,
	// the apns-topic header field must use your app’s bundle ID with .pushkit.fileprovider appended to the end.
	PushTypeFileprovider = PushType("fileprovider")

	// Use the mdm push type for notifications that tell managed devices to contact the MDM server. If you set
	// this push type, you must use the topic from the UID attribute in the subject of your MDM push certificate.
	PushTypeMdm = PushType("mdm")
)

type Headers struct {
	// (Required for watchOS 6 and later; recommended for macOS, iOS, tvOS, and iPadOS) The value of this header
	// must accurately reflect the contents of your notification’s payload. If there is a mismatch, or if the
	// header is missing on required systems, APNs may return an error, delay the delivery of the notification,
	// or drop it altogether.
	PushType PushType

	// A canonical UUID that is the unique ID for the notification. If an error occurs when sending the notification,
	// APNs includes this value when reporting the error to your server. Canonical UUIDs are 32 lowercase hexadecimal
	// digits, displayed in five groups separated by hyphens in the form 8-4-4-4-12. An example looks like this:
	// 123e4567-e89b-12d3-a456-4266554400a0. If you omit this header, APNs creates a UUID for you and returns it
	// in its response.
	Id string

	// The date at which the notification is no longer valid. This value is a UNIX epoch expressed in seconds (UTC).
	// If the value is nonzero, APNs stores the notification and tries to deliver it at least once, repeating the
	// attempt as needed until the specified date. If the value is 0, APNs attempts to deliver the notification
	// only once and doesn’t store it.
	Expiration time.Time

	// The priority of the notification. If you omit this header, APNs sets the notification priority to 10.
	// Specify 10 to send the notification immediately. A value of 10 is appropriate for notifications that trigger
	// an alert, play a sound, or badge the app’s icon. Specifying this priority for a notification that has a payload
	// containing the content-available key causes an error. Specify 5 to send the notification based on power
	// considerations on the user’s device. Use this priority for notifications that have a payload that includes
	// the content-available key. Notifications with this priority might be grouped and delivered in bursts to the
	// user’s device. They may also be throttled, and in some cases not delivered.
	Priority int

	// An identifier you use to coalesce multiple notifications into a single notification for the user. Typically,
	// each notification request causes a new notification to be displayed on the user’s device. When sending the
	// same notification more than once, use the same value in this header to coalesce the requests. The value
	// of this key must not exceed 64 bytes.
	CollapseId string

	// The topic for the notification. In general, the topic is your app’s bundle ID, but it may have a suffix
	// based on the push notification’s type.
	topic string
}

func (h Headers) Map() map[string]string {
	res := make(map[string]string)

	if h.PushType != "" {
		res["apns-push-type"] = string(h.PushType)
	}

	if h.Id != "" {
		res["apns-id"] = h.Id
	}

	if !h.Expiration.IsZero() {
		res["apns-expiration"] = fmt.Sprintf("%d", h.Expiration.Unix())
	}

	if h.Priority > 0 {
		res["apns-priority"] = fmt.Sprintf("%d", h.Priority)
	}

	if h.topic != "" {
		res["apns-topic"] = h.topic
		switch h.PushType {
		case PushTypeVoip:
			res["apns-topic"] += ".voip"
		case PushTypeComplication:
			res["apns-topic"] += ".complication"
		case PushTypeFileprovider:
			res["apns-topic"] += ".pushkit.fileprovider"
		}
	}

	return res
}
