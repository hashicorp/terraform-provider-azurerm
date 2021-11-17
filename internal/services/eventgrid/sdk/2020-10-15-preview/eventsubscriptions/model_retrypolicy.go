package eventsubscriptions

type RetryPolicy struct {
	EventTimeToLiveInMinutes *int64 `json:"eventTimeToLiveInMinutes,omitempty"`
	MaxDeliveryAttempts      *int64 `json:"maxDeliveryAttempts,omitempty"`
}
