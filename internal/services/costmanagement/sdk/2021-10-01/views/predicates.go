package views

type ViewPredicate struct {
	ETag *string
	Id   *string
	Name *string
	Type *string
}

func (p ViewPredicate) Matches(input View) bool {

	if p.ETag != nil && (input.ETag == nil && *p.ETag != *input.ETag) {
		return false
	}

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
