package locations

type Trial struct {
	AvailableHosts *int64       `json:"availableHosts,omitempty"`
	Status         *TrialStatus `json:"status,omitempty"`
}
