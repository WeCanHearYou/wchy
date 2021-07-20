package enum

// WebhookStatus is the status of a webhook
type WebhookStatus int

const (
	// WebhookEnabled means the webhook can be triggered
	WebhookEnabled WebhookStatus = 1
	// WebhookDisabled means the webhook cannot be triggered
	WebhookDisabled WebhookStatus = 2
	// WebhookFailed means an error occured when the webhook was previously triggered and has been disabled
	WebhookFailed WebhookStatus = 3
)

var webhookStatusIDs = map[WebhookStatus]string{
	WebhookEnabled:  "enabled",
	WebhookDisabled: "disabled",
	WebhookFailed:   "failed",
}

var webhookStatusName = map[string]WebhookStatus{
	"enabled":  WebhookEnabled,
	"disabled": WebhookDisabled,
	"failed":   WebhookFailed,
}

// MarshalText returns the Text version of the webhook status
func (status WebhookStatus) MarshalText() ([]byte, error) {
	return []byte(webhookStatusIDs[status]), nil
}

// UnmarshalText parse string into a webhook status
func (status *WebhookStatus) UnmarshalText(text []byte) error {
	*status = webhookStatusName[string(text)]
	return nil
}

// Name returns the name of a webhook status
func (status WebhookStatus) Name() string {
	name, ok := webhookStatusIDs[status]
	if ok {
		return name
	}
	return "unknown"
}
