package capacities

type CheckCapacityNameAvailabilityResult struct {
	Message       *string `json:"message,omitempty"`
	NameAvailable *bool   `json:"nameAvailable,omitempty"`
	Reason        *string `json:"reason,omitempty"`
}
