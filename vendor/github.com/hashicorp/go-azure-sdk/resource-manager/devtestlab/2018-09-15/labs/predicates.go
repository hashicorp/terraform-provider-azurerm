package labs

type LabOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p LabOperationPredicate) Matches(input Lab) bool {

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

type LabVhdOperationPredicate struct {
	Id *string
}

func (p LabVhdOperationPredicate) Matches(input LabVhd) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	return true
}
