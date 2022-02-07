package videoanalyzer

type CheckNameAvailabilityResponse struct {
	Message       *string                      `json:"message,omitempty"`
	NameAvailable *bool                        `json:"nameAvailable,omitempty"`
	Reason        *CheckNameAvailabilityReason `json:"reason,omitempty"`
}
