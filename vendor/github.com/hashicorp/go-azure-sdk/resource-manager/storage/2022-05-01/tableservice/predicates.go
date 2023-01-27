package tableservice

type TableOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p TableOperationPredicate) Matches(input Table) bool {

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
