package checknameavailabilitydisasterrecoveryconfigs

type CheckNameAvailabilityResult struct {
	Message       *string            `json:"message,omitempty"`
	NameAvailable *bool              `json:"nameAvailable,omitempty"`
	Reason        *UnavailableReason `json:"reason,omitempty"`
}
