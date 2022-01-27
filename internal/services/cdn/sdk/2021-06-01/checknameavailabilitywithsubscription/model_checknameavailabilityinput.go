package checknameavailabilitywithsubscription

type CheckNameAvailabilityInput struct {
	Name string       `json:"name"`
	Type ResourceType `json:"type"`
}
