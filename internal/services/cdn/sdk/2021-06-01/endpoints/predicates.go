package endpoints

type EndpointPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p EndpointPredicate) Matches(input Endpoint) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && *p.Location != input.Location {
		return false
	}

	if p.Name != nil && (input.Name == nil && *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil && *p.Type != *input.Type) {
		return false
	}

	return true
}

type ResourceUsagePredicate struct {
	CurrentValue *int64
	Limit        *int64
	ResourceType *string
	Unit         *string
}

func (p ResourceUsagePredicate) Matches(input ResourceUsage) bool {

	if p.CurrentValue != nil && (input.CurrentValue == nil && *p.CurrentValue != *input.CurrentValue) {
		return false
	}

	if p.Limit != nil && (input.Limit == nil && *p.Limit != *input.Limit) {
		return false
	}

	if p.ResourceType != nil && (input.ResourceType == nil && *p.ResourceType != *input.ResourceType) {
		return false
	}

	if p.Unit != nil && (input.Unit == nil && *p.Unit != *input.Unit) {
		return false
	}

	return true
}
