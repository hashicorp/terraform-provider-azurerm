package get

type GroupIdInformationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p GroupIdInformationPredicate) Matches(input GroupIdInformation) bool {

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

type PrivateEndpointConnectionPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p PrivateEndpointConnectionPredicate) Matches(input PrivateEndpointConnection) bool {

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
