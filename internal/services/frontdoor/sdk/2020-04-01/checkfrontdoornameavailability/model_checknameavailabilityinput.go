package checkfrontdoornameavailability

type CheckNameAvailabilityInput struct {
	Name string       `json:"name"`
	Type ResourceType `json:"type"`
}
