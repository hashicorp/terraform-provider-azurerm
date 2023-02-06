package linkedresources

type LinkedResourceOperationPredicate struct {
	Id *string
}

func (p LinkedResourceOperationPredicate) Matches(input LinkedResource) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	return true
}
