package signalr

type PrivateLinkResourcePredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p PrivateLinkResourcePredicate) Matches(input PrivateLinkResource) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
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

type SignalRResourcePredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p SignalRResourcePredicate) Matches(input SignalRResource) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil && *p.Location != *input.Location) {
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

type SignalRUsagePredicate struct {
	CurrentValue *int64
	Id           *string
	Limit        *int64
	Unit         *string
}

func (p SignalRUsagePredicate) Matches(input SignalRUsage) bool {

	if p.CurrentValue != nil && (input.CurrentValue == nil && *p.CurrentValue != *input.CurrentValue) {
		return false
	}

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	if p.Limit != nil && (input.Limit == nil && *p.Limit != *input.Limit) {
		return false
	}

	if p.Unit != nil && (input.Unit == nil && *p.Unit != *input.Unit) {
		return false
	}

	return true
}
