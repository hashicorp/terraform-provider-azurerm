package accounts

type CheckNameAvailabilityParameters struct {
	Name string `json:"name"`
	Type Type   `json:"type"`
}
