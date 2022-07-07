package videoanalyzer

type AccessPolicyEntityOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p AccessPolicyEntityOperationPredicate) Matches(input AccessPolicyEntity) bool {

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

type EdgeModuleEntityOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p EdgeModuleEntityOperationPredicate) Matches(input EdgeModuleEntity) bool {

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

type VideoEntityOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p VideoEntityOperationPredicate) Matches(input VideoEntity) bool {

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
