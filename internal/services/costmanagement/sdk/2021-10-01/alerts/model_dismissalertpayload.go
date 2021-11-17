package alerts

type DismissAlertPayload struct {
	Properties *AlertProperties `json:"properties,omitempty"`
}
