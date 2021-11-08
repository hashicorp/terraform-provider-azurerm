package frontdoors

type FrontDoorPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p FrontDoorPredicate) Matches(input FrontDoor) bool {

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

type FrontendEndpointPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p FrontendEndpointPredicate) Matches(input FrontendEndpoint) bool {

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

type RulesEnginePredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RulesEnginePredicate) Matches(input RulesEngine) bool {

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
