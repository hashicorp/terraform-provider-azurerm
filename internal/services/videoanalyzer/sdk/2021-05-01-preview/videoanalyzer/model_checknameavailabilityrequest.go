package videoanalyzer

type CheckNameAvailabilityRequest struct {
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}
