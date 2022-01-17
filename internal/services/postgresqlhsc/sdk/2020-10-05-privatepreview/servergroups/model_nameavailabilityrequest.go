package servergroups

type NameAvailabilityRequest struct {
	Name string                            `json:"name"`
	Type CheckNameAvailabilityResourceType `json:"type"`
}
