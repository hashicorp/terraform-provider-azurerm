package checkfrontdoornameavailabilitywithsubscription

type CheckNameAvailabilityOutput struct {
	Message          *string       `json:"message,omitempty"`
	NameAvailability *Availability `json:"nameAvailability,omitempty"`
	Reason           *string       `json:"reason,omitempty"`
}
