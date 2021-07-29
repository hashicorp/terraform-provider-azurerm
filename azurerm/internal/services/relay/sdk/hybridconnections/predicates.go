package hybridconnections

type AuthorizationRulePredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p AuthorizationRulePredicate) Matches(input AuthorizationRule) bool {

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

type HybridConnectionPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p HybridConnectionPredicate) Matches(input HybridConnection) bool {

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
